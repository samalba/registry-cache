vcl 4.1;

import directors;

# Define a shared probe
probe docker_probe {
    .url = "/v2/";
    .timeout = 5s;
    .interval = 10s;
    .window = 5;
    .threshold = 3;
}

# Define multiple backends for redundancy
backend docker_registry_backend_1 {
    .host = "54.198.86.24";
    .port = "443";
    .probe = docker_probe;  # Use the shared probe
}

backend docker_registry_backend_2 {
    .host = "54.236.113.205";
    .port = "443";
    .probe = docker_probe;  # Use the shared probe
}

backend docker_registry_backend_3 {
    .host = "54.227.20.253";
    .port = "443";
    .probe = docker_probe;  # Use the shared probe
}

# Initialize the round-robin director
sub vcl_init {
    new docker_director = directors.round_robin();
    docker_director.add_backend(docker_registry_backend_1);
    docker_director.add_backend(docker_registry_backend_2);
    docker_director.add_backend(docker_registry_backend_3);
}

sub vcl_recv {
    # Use the round-robin director for load balancing
    set req.backend_hint = docker_director.backend();

    # Set Varnish to operate as a general HTTP proxy
    if (req.method == "GET" || req.method == "HEAD") {
        # Handle OCI image layer caching logic

        # Check for OCI image layer URL patterns (assuming the layer URLs contain "/blobs/sha256")
        if (req.url ~ "^/v2/.*/blobs/sha256") {
            # Allow caching of image layers
            return (hash);
        } else if (req.url ~ "^/v2/.*/manifests/") {
            # Bypass cache for manifest files
            return (pass);
        }
    }

    # For non-GET/HEAD requests or other traffic, pass through the proxy without caching
    return (pass);
}

sub vcl_backend_response {
    # Cache the image layers without setting a TTL
    if (bereq.url ~ "^/v2/.*/blobs/sha256") {
        # Cache the image layers
        set beresp.ttl = 24h;  # Set a default TTL if needed
    } else {
        # Do not cache other responses
        set beresp.ttl = 0s;  # Effectively disable caching for non-layer responses
        set beresp.uncacheable = true;
    }
}

sub vcl_deliver {
    # Optionally, add headers to show whether the response was served from cache
    if (obj.hits > 0) {
        set resp.http.X-Cache = "HIT";
    } else {
        set resp.http.X-Cache = "MISS";
    }
}
