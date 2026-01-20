package dom

import "math/rand"

const TeleportBlock = "__TELEPORT_BLOCK"

const PageServerProps = "__PAGE_PROPS"

const (
	ToastContainerTopRight  = "__TOAST_CONTAINER_TOP_RIGHT"
	ToastContainerTopLeft   = "__TOAST_CONTAINER_TOP_LEFT"
	ToastContainerTopCenter = "__TOAST_CONTAINER_TOP_CENTER"

	ToastContainerBottomRight  = "__TOAST_CONTAINER_BOTTOM_RIGHT"
	ToastContainerBottomLeft   = "__TOAST_CONTAINER_BOTTOM_LEFT"
	ToastContainerBottomCenter = "__TOAST_CONTAINER_BOTTOM_CENTER"
)

var toastContainers = []string{
	ToastContainerTopRight,
	ToastContainerTopLeft,
	ToastContainerTopCenter,
	ToastContainerBottomRight,
	ToastContainerBottomLeft,
	ToastContainerBottomCenter,
}

// RandomToastContainerID возвращает случайный id toast-контейнера
func RandomToastContainerID() string {
	return toastContainers[rand.Intn(len(toastContainers))]
}

const META_CSRF = "meta-csrf"
