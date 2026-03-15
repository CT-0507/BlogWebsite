let worker: SharedWorker | null = null;

function getWorker() {
  if (!worker) {
    worker = new SharedWorker("/sseWorker.js");
    worker.port.start();
  }

  return worker;
}

export function initAuthSSE(
  baseURL: string,
  token: string,
  topics: string[],
  globalTopics?: string[]
) {
  const w = getWorker();

  w.port.postMessage({
    type: "init-auth",
    baseURL,
    token,
    topics,
    globalTopics,
  });

  console.log(w);

  return w;
}

export function initPublicSSE(
  baseURL: string,
  topics?: string[],
  globalTopics?: string[]
) {
  const w = getWorker();

  w.port.postMessage({
    type: "init-public",
    baseURL,
    topics,
    globalTopics,
  });

  return w;
}
