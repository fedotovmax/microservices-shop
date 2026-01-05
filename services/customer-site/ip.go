func getClientIP(r *http.Request) string {
	// X-Forwarded-For от внешнего прокси (если есть)
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		parts := strings.Split(xff, ",")
		return strings.TrimSpace(parts[0])
	}

	ip, _, _ := net.SplitHostPort(r.RemoteAddr)

	if ip == "::1" {
		return "127.0.0.1"
	}

	parsed := net.ParseIP(ip)
	if v4 := parsed.To4(); v4 != nil {
		return v4.String()
	}

	return ip
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	clientIP := getClientIP(r)
	userAgent := r.UserAgent()

	body := map[string]string{
		"email":    r.FormValue("email"),
		"password": r.FormValue("password"),
	}

	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest(
		http.MethodPost,
		"http://api-gateway/login",
		bytes.NewBuffer(jsonBody),
	)

	req.Header.Set("Content-Type", "application/json")

	// ✅ Прокидываем User-Agent
	req.Header.Set("User-Agent", userAgent)

	// ✅ Прокидываем IP
	req.Header.Set("X-Real-IP", clientIP)

	// Если Go-сервер сам является прокси
	if prior := r.Header.Get("X-Forwarded-For"); prior != "" {
		req.Header.Set("X-Forwarded-For", prior+", "+clientIP)
	} else {
		req.Header.Set("X-Forwarded-For", clientIP)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		http.Error(w, "Gateway error", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}
