import {
  forwardRef,
  memo,
  useEffect,
  useImperativeHandle,
  useRef,
  useState,
} from "react";
import { type OutputData } from "@editorjs/editorjs";
import Box from "@mui/material/Box";
import type EditorJS from "@editorjs/editorjs";
import { uploadByFile } from "@/api/blogApi";
import debounce from "lodash.debounce";
import Button from "@mui/material/Button";
import Dialog from "@mui/material/Dialog";
import DialogActions from "@mui/material/DialogActions";
import DialogContent from "@mui/material/DialogContent";
import DialogContentText from "@mui/material/DialogContentText";
import DialogTitle from "@mui/material/DialogTitle";
import "./editor.css";
import { RECOVERY_KEY } from "./BlogForm";

export interface EditorSavedData {
  content: OutputData | null;

  files: Map<string, File>;
}

export interface EditorHandle {
  save: () => Promise<EditorSavedData>;
}

interface EditorProps {
  initialData?: OutputData;
}

interface RecoverDialogProps {
  handleRecover: () => Promise<void>;
  handleDiscard: () => void;
}
export function RecoverDialog({
  handleRecover,
  handleDiscard,
}: RecoverDialogProps) {
  const [open, setOpen] = useState(true);
  const [openDiscardConfirm, setOpenDiscardConfirm] = useState(false);
  const handleDiscardClick = () => {
    setOpenDiscardConfirm(true);
  };
  const handleRecoverClick = async () => {
    await handleRecover();
    setOpen(false);
  };
  const handleCloseDiscardConfirmDialog = () => {
    setOpenDiscardConfirm(false);
  };
  const handleDiscardDialogConfirm = () => {
    handleDiscard();
    setOpenDiscardConfirm(false);
    setOpen(false);
  };
  const handleOpenConfirmDialog = () => {
    setOpenDiscardConfirm(true);
  };
  return (
    <>
      <Dialog
        open={open}
        onClose={handleDiscardClick}
        aria-labelledby="alert-dialog-title"
        aria-describedby="alert-dialog-description"
        role="alertdialog"
      >
        <DialogTitle id="alert-dialog-title">{"Content recover?"}</DialogTitle>
        <DialogContent>
          <DialogContentText id="alert-dialog-description">
            We discover that you have unsaved content before. Do you want to
            recover it?
          </DialogContentText>
        </DialogContent>
        <DialogActions>
          <Button onClick={handleOpenConfirmDialog} color="error">
            Discard
          </Button>
          <Button onClick={handleRecoverClick} autoFocus>
            Recover
          </Button>
        </DialogActions>
      </Dialog>
      <Dialog
        open={openDiscardConfirm}
        onClose={handleCloseDiscardConfirmDialog}
        aria-labelledby="alert-dialog-title"
        aria-describedby="alert-dialog-description"
        role="alertdialog"
      >
        <DialogTitle id="alert-dialog-title">
          {"You want to discard the content?"}
        </DialogTitle>
        <DialogContent></DialogContent>
        <DialogActions>
          <Button onClick={handleDiscardDialogConfirm} color="error">
            Confirm
          </Button>
          <Button onClick={handleCloseDiscardConfirmDialog} autoFocus>
            Cancel
          </Button>
        </DialogActions>
      </Dialog>
    </>
  );
}

const Editor = forwardRef<EditorHandle, EditorProps>(({ initialData }, ref) => {
  // Recovery
  // let recovery = localStorage.getItem(RECOVERY_KEY);

  const holderRef = useRef<HTMLDivElement | null>(null);
  const editorRef = useRef<EditorJS | null>(null);
  const fileMapRef = useRef<Map<string, File>>(new Map());

  useImperativeHandle(ref, () => ({
    async save() {
      if (!editorRef.current) {
        return {
          content: null,
          files: new Map(),
        };
      }

      const content = await editorRef.current.save();

      return {
        content,

        files: fileMapRef.current,
      };
    },
  }));

  const saveRecovery = debounce(async () => {
    if (!editorRef.current) return;

    const content = await editorRef.current.save();

    localStorage.setItem(
      RECOVERY_KEY,
      JSON.stringify({
        content,
        updatedAt: Date.now(),
      }),
    );
  }, 2000);
  useEffect(() => {
    let editor: EditorJS;
    if (!holderRef.current) return;
    const init = async () => {
      const EditorJSB = (await import("@editorjs/editorjs")).default;
      const ImageTool = (await import("@editorjs/image")).default;
      const List = (await import("@editorjs/list")).default;
      const Header = (await import("@editorjs/header")).default;
      const LinkTool = (await import("@editorjs/link")).default;
      if (!editorRef.current) {
        editor = new EditorJSB({
          holder: holderRef.current!,
          async onChange() {
            saveRecovery();
          },
          inlineToolbar: true,
          data: initialData ?? { blocks: [] },
          tools: {
            list: {
              class: List,
              inlineToolbar: true,
              config: {
                defaultStyle: "unordered",
              },
            },
            header: {
              // eslint-disable-next-line @typescript-eslint/no-explicit-any
              class: Header as any,
              inlineToolbar: true,
            },
            image: {
              class: ImageTool,
              config: {
                uploader: {
                  uploadByFile: uploadByFile,
                },
              },
            },
            link: {
              class: LinkTool,
              shortcut: "CMD+SHIFT+H",
              inlineToolbar: true,
            },
          },
          placeholder: "Write something",
        });

        editorRef.current = editor;
      }
    };

    if (typeof window !== "undefined") {
      init().then();
    }
    return () => {
      editor?.isReady?.then(() => editor.destroy()).catch(() => null);
    };
  }, [initialData, editorRef, saveRecovery]);

  // const handleRecoverContent = async () => {
  //   if (editorRef && editorRef.current && recovery) {
  //     const parsed = JSON.parse(recovery);
  //     editorRef.current.render(parsed);
  //   }
  // };

  // const handleDiscard = async () => {
  //   recovery = null;
  //   localStorage.removeItem(RECOVERY_KEY);
  // };

  return (
    <Box sx={{ p: 0, m: 0, width: "100%" }}>
      {/* <Box>
        <Button onClick={handleRecoverContent} disabled={!recovery}>
          Recovery last content
        </Button>
      </Box> */}
      {/* {!!recovery && (
        <RecoverDialog
          handleDiscard={handleDiscard}
          handleRecover={handleRecoverContent}
        />
      )} */}
      <Box
        sx={{ border: "2px solid gray", borderRadius: 2, p: 0 }}
        ref={holderRef}
      />
    </Box>
  );
});

export default memo(Editor);
