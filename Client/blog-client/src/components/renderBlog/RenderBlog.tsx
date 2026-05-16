import Typography from "@mui/material/Typography";
import DOMPurify from "dompurify";

type Block = {
  type: string;
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  data: any;
};

type EditorContent = {
  blocks: Block[];
};

type HeaderVariant = "h1" | "h2" | "h3" | "h4";

export function RenderArticle({ content }: { content: EditorContent }) {
  return (
    <article className="article">
      {content.blocks.map((block, index) => {
        switch (block.type) {
          case "header": {
            const level = block.data.level;
            const text = block.data.text;
            let variant: HeaderVariant = "h4";
            switch (level) {
              case 1:
                variant = "h1";
                break;
              case 2:
                variant = "h2";
                break;

              case 3:
                variant = "h3";
                break;

              default:
                variant = "h4";
                break;
            }
            return (
              <Typography key={index} variant={variant}>
                {text}
              </Typography>
            );
          }

          case "paragraph":
            return (
              <p
                key={index}
                dangerouslySetInnerHTML={{
                  __html: DOMPurify.sanitize(block.data.text),
                }}
              />
            );

          case "list": {
            const ListTag = block.data.style === "ordered" ? "ol" : "ul";

            return (
              <ListTag key={index}>
                {block.data.items.map((item: string, i: number) => (
                  <li
                    key={i}
                    dangerouslySetInnerHTML={{
                      __html: DOMPurify.sanitize(item),
                    }}
                  />
                ))}
              </ListTag>
            );
          }

          case "image":
            return (
              <figure key={index}>
                <img src={block.data.file.url} alt={block.data.caption || ""} />

                {block.data.caption && (
                  <figcaption>{block.data.caption}</figcaption>
                )}
              </figure>
            );

          case "linkTool":
            return (
              <a
                key={index}
                href={block.data.link}
                target="_blank"
                rel="noopener noreferrer"
              >
                {block.data.meta?.title || block.data.link}
              </a>
            );

          default:
            return null;
        }
      })}
    </article>
  );
}
