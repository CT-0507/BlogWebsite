/* eslint-disable @typescript-eslint/no-explicit-any */
import { useEffect } from "react";
import { useQueryClient } from "@tanstack/react-query";
import { initAuthSSE, initPublicSSE } from "../lib/sseClient";
import { BASE_URL } from "@/api/axiosConfig";

export function useAuthSSE(
  token: string | null,
  topics: string[],
  globalTopics?: string[],
  setSnackbar?: (value: boolean) => void
) {
  const queryClient = useQueryClient();

  useEffect(() => {
    if (!token) return;

    const worker = initAuthSSE(BASE_URL, token, topics, globalTopics);

    worker.port.start();

    worker.port.onmessage = (msg) => {
      console.log(msg);
      if (msg.data.type !== "cache-patch") return;

      const { queryKey, op, data } = msg.data.patch;

      setSnackbar?.(true);

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
  }, [topics, token, queryClient, globalTopics, setSnackbar]);
}

export function usePublicSSE(topics: string[], globalTopics?: string[]) {
  const queryClient = useQueryClient();

  useEffect(() => {
    const worker = initPublicSSE(BASE_URL, topics, globalTopics);

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
