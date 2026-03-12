let worker: SharedWorker | null = null;

function getWorker() {
  if (!worker) {
    worker = new SharedWorker("/sseWorker.js");
    worker.port.start();
  }

  return worker;
}

export function initAuthSSE(
  token: string,
  topics?: string[],
  globalTopics?: string[]
) {
  console.log("Auth w");
  const w = getWorker();

  w.port.postMessage({
    type: "init-auth",
    token,
    topics,
    globalTopics,
  });

  console.log(w);

  return w;
}

export function initPublicSSE(topics?: string[], globalTopics?: string[]) {
  const w = getWorker();

  w.port.postMessage({
    type: "init-public",
    topics,
    globalTopics,
  });

  return w;
}
