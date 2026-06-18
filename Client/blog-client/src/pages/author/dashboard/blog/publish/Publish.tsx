import { appName } from "@/config/const";
import BlogForm from "../components/BlogForm";

export default function Publish() {
  const title = `${appName} | Publish`;
  return (
    <>
      <title>{title}</title>
      <BlogForm mode="create" />
    </>
  );
}
