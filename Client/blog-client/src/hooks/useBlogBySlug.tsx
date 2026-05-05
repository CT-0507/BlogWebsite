import { getBlogBySlug } from "@/api/blogApi";
import { useQuery } from "@tanstack/react-query";

export function useBlogBySlug(isAuthenticated: boolean, slug?: string) {
  // const queryClient = useQueryClient();

  return useQuery({
    queryKey: ["blog", slug],
    queryFn: () => getBlogBySlug(slug!, isAuthenticated),
    staleTime: Infinity,
    enabled: !!slug,

    // initialData: () => {
    //   // When using pagination
    //   //   const pages = queryClient.getQueryData(["items"]);

    //   //   if (!pages) return undefined;

    //   //   for (const page of pages.pages) {
    //   //     const found = page.items.find((i) => i.id === itemId);
    //   //     if (found) return found;
    //   //   }
    //   const blogs = queryClient.getQueryData<Blog[]>(["blogs"]);
    //   return blogs?.find((i) => i.urlSlug === slug);
    // },
  });
}
