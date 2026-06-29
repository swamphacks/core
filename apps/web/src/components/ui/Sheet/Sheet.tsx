import {
  Modal,
  ModalOverlay,
  type ModalOverlayProps,
  Heading,
  composeRenderProps,
} from "react-aria-components";
import { Dialog } from "../Dialog";
import "./Sheet.css";
import { cn } from "@/utils/cn";

type SheetProps = ModalOverlayProps & {
  sheetClassName?: string;
};

export function Sheet(props: SheetProps) {
  return (
    <ModalOverlay isDismissable className="sheet-overlay">
      {composeRenderProps(props.children, (children) => (
        <Modal className={cn("sheet", props.sheetClassName)}>
          <Dialog>{children}</Dialog>
        </Modal>
      ))}
    </ModalOverlay>
  );
}

export { Heading };
