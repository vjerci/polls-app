const API_URL = process.env.NEXT_PUBLIC_API_URL;

const doRequest =
	(method: string) =>
	async ({
		path,
		headers,
		body,
	}: {
		path: string;
		headers?: any;
		body?: any;
	}) => {
		const config: RequestInit = {};

		if (typeof body !== "undefined") {
			config.body = JSON.stringify(body);
		}

		config.headers = {};
		config.headers["Content-Type"] = "application/json";

		if (typeof headers !== "undefined") {
			config.headers = { ...config.headers, ...headers };
		}

		config.method = method;

		const resp = await fetch(`${API_URL}${path}`, config);

		return await resp.json();
	};

export const api = {
	POST: doRequest("POST"),
	PUT: doRequest("PUT"),
	GET: doRequest("GET"),
	DELETE: doRequest("DELETE"),
};

export const buildAuthHeaders = (token: string) => {
	return {
		Authorization: `Bearer ${token}`,
	};
};

export default api;
