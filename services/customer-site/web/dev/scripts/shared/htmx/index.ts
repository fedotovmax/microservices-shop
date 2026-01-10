import htmx from 'htmx.org';

type HtmxValues = Record<string, any>;
type HtmxHeaders = Record<string, string>;

type HttpVerb =
	| 'get'
	| 'head'
	| 'post'
	| 'put'
	| 'delete'
	| 'connect'
	| 'options'
	| 'trace'
	| 'patch';

export function htmxRequest(
	url: string,
	method: HttpVerb,
	data: HtmxValues = {},
	headers: HtmxHeaders = {}
): Promise<string> {
	return new Promise((resolve, reject) => {
		const handler = (evt: Event) => {
			const e = evt as CustomEvent;
			if (e.detail.pathInfo.requestPath === url) {
				document.body.removeEventListener('htmx:afterRequest', handler);
				const status = e.detail.xhr.status;
				if (status >= 200 && status < 300) {
					resolve('ok');
				} else {
					reject(new Error(`HTMX request failed with status ${status}`));
				}
			}
		};

		document.body.addEventListener('htmx:afterRequest', handler);

		htmx.ajax(method, url, {
			values: data,
			headers: headers,
		});
	});
}
