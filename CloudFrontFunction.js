function handler(event) {
    var request = event.request;
    var headers = request.headers;

    if (!headers.host || !headers.host.value) {
        return {
            statusCode: 400,
            statusDescription: "Bad Request",
            headers: {
                "content-type": { value: "text/plain" }
            },
            body: "Host header missing"
        };
    }

    var host = headers.host.value;
    var subdomain = host.split('.')[0];
    var uri = request.uri;

    // List of static file extensions
    var staticExts = [
        ".js", ".css", ".png", ".jpg", ".jpeg", ".gif", ".ico", ".svg",
        ".woff", ".woff2", ".ttf", ".eot", ".map", ".json"
    ];

    // Helper to check if URI is a static file
    function isStaticFile(path) {
        return staticExts.some(function(ext) {
            return path.endsWith(ext);
        });
    }

    if (uri === "/") {
        request.uri = "/" + subdomain + "/index.html";
    } else if (isStaticFile(uri)) {
        request.uri = "/" + subdomain + uri;
    } else {
        // For SPA routes, always serve index.html
        request.uri = "/" + subdomain + "/index.html";
    }

    return request;
}
