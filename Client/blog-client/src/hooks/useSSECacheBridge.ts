/* eslint-disable @typescript-eslint/no-explicit-any */
import { useEffect } from "react";
import { useQueryClient } from "@tanstack/react-query";
import { initAuthSSE, initPublicSSE } from "../lib/sseClient";

export function useAuthSSE(
  token: string | null,
  topics?: string[],
  globalTopics?: string[]
) {
  console.log("Run");
  const queryClient = useQueryClient();

  useEffect(() => {
    if (!token) return;

    console.log(globalTopics);

    fetch("http://localhost:8080/events/auth?topics=blog_created_admin", {
      headers: {
        Accept: "text/event-stream",
        Authorization: `Bearer ${token}`,
      },
    }).then(async (res) => {
      const reader = res.body.getReader();
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

          const payload = JSON.parse(line.replace("Cache:", "").trim());

          console.log(payload);
        }
      }
    });

    const worker = initAuthSSE(token, topics, globalTopics);

    worker.port.start();

    worker.port.onmessage = (msg) => {
      console.log(msg);
      if (msg.data.type !== "cache-patch") return;

      const { queryKey, op, data } = msg.data.patch;

      queryClient.setQueryData(queryKey, (old: any) => {
        if (!old) return old;

        switch (op) {
          case "append":
            return [...old, data];

          case "prepend":
            return [data, ...old];

          case "merge":
            return { ...old, ...data };

          case "remove":
            return old.filter((x: any) => x.id !== data.id);

          case "replace":
            return data;

          default:
            return old;
        }
      });
    };
  }, [topics, token, queryClient, globalTopics]);
}

export function usePublicSSE(topics: string[], globalTopics?: string[]) {
  const queryClient = useQueryClient();

  useEffect(() => {
    const worker = initPublicSSE(topics, globalTopics);

    worker.port.onmessage = (msg) => {
      if (msg.data.type !== "cache-patch") return;

      const { queryKey, op, data } = msg.data.patch;

      queryClient.setQueryData(queryKey, (old: any) => {
        if (!old) return old;

        switch (op) {
          case "append":
            return [...old, data];

          case "merge":
            return { ...old, ...data };

          default:
            return old;
        }
      });
    };
  }, [queryClient, globalTopics, topics]);
}
