const ports = new Map();

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

let retryCount = 0;
const MAX_RETRIES = 10;
const BASE_DELAY = 1000;
const BASE_URL = "http://localhost:8080";

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
      const line = part.split("\n").find((l) => l.startsWith("content:"));
      if (!line) continue;

      const payload = JSON.parse(line.replace("content:", "").trim());

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
  for (var topic of topics) {
    queryParams += topic;
  }
  console.log(queryParams);
  return queryParams;
}

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

async function startPublicStream(topics) {
  if (publicController) return;

  publicController = new AbortController();

  await startStream("/events/public" + mapTopicsToQueryParams(topics), publicController, (event) =>
    routeEvent(event, "public")
  );
}

function routeEvent(event, type) {

  if (globalTopics.has(event.topic)) {

    // broadcast to all tabs
    for (const port of portTopics.keys()) {

      port.postMessage({
        type: "cache-patch",
        patch: event.cache
      })
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
    const topics = msg.topics

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
