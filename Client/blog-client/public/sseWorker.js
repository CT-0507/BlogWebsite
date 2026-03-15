const ports = new Map();

// Store topics that are broadcasted to all tabs
const globalTopics = new Set();

// Expected message for set cache
// {
//   "topic": "comments:123",
//   "cache": {
//     "queryKey": ["comments", 123],
//     "op": "append",
//     "data": { "id": 9, "text": "hello" }
//   }
// }

let authController = null;
let publicController = null;

let token = null;

// TODO: auto retry on server connection failure
let retryCount = 0;
const MAX_RETRIES = 10;
const BASE_DELAY = 1000;

// SSE server URL
let BASE_URL;

async function startStream(
  url,
  controller,
  route,
  headers
) {
  const response = await fetch(BASE_URL + url, {
    headers: {
      Accept: "text/event-stream",
      ...headers,
    },
    signal: controller.signal,
  });

  // Manually decode message from server
  const reader = response.body.getReader();
  const decoder = new TextDecoder();

  let buffer = "";

  while (true) {

    const { done, value } = await reader.read();

    if (done) break;

    buffer += decoder.decode(value, { stream: true });

    const parts = buffer.split("\n\n");
    buffer = parts.pop() || "";

    for (const part of parts) {

      // Take only data property to broadcast
      const line = part.split("\n").find((l) => l.startsWith("data:"));
      if (!line) continue;

      const payload = JSON.parse(line.replace("data:", "").trim());

      route(payload);
    }
  }
}

function addGlobalTopics(topics) {

  for (const topic of topics) {

    if (globalTopics.has(topic)) continue

    globalTopics.add(topic)
  }
}

function mapTopicsToQueryParams(topics) {
  let queryParams = "?topics=";
  for (var topic of globalTopics) {
    queryParams += topic + ",";
  }
  for (var topic of topics) {
    queryParams += topic + ",";
  }
  queryParams = queryParams.slice(0, -1);
  return queryParams;
}

/** Start authorized SSE */ 
async function startAuthStream(topics) {
  if (authController || !token) return;

  authController = new AbortController();

  await startStream(
    "/events/auth" + mapTopicsToQueryParams(topics),
    authController,
    (event) => routeEvent(event, "auth"),
    { Authorization: `Bearer ${token}` }
  );
}

/** Start public SSE */ 
async function startPublicStream(topics) {
  if (publicController) return;

  publicController = new AbortController();

  await startStream("/events/public" + mapTopicsToQueryParams(topics), publicController, (event) =>
    routeEvent(event, "public")
  );
}

/** Route event to its corresponding tab and route global event to all tab */
function routeEvent(event, type) {

  if (globalTopics.has(event.topic)) {

    // broadcast to all tabs
    for (const port of portTopics.keys()) {

      // Only update cache event for now
      port.postMessage({
        type: "cache-patch",
        patch: event.cache
      })

      // TODO: implement other kind of messages
    }

    return
  }

  for (const [port, meta] of ports.entries()) {
    if (type === "auth" && !meta.auth) continue;
    if (type === "public" && !meta.public) continue;

    port.postMessage({
      type: "cache-patch",
      patch: event.cache,
    });
  }
}

onconnect = (e) => {

  const port = e.ports[0]

  ports.set(port, { auth: false, public: false })

  port.start()

  port.onmessage = msg => {

    const data = msg.data
    const topics = data.topics
    BASE_URL = data.baseURL

    if (data.type === "init-auth") {

      token = data.token

      const meta = ports.get(port)
      meta.auth = true

      if (data.globalTopics)
        addGlobalTopics(data.globalTopics)

      startAuthStream(topics)
    }

    if (data.type === "init-public") {

      const meta = ports.get(port)
      meta.public = true

      startPublicStream(topics)
    }
  }

  port.onmessageerror = () => {
    ports.delete(port)
  }
}
