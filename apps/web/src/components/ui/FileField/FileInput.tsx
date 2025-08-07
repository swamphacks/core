import { type JSX, useContext, useRef } from "react";
import {
  DropZone,
  InputContext,
  useContextProps,
  type FileDropItem,
  type InputProps,
} from "react-aria-components";
import { type ChangeEvent, useCallback, useState } from "react";
import { Input } from "@/components/ui/Field";
import { cn } from "@/utils/cn";
import { inputStyles } from "@/components/ui/TextField";
import { Button } from "@/components/ui/Button";
import TablerX from "~icons/tabler/x";
import TablerUpload from "~icons/tabler/upload";
import { CustomPropsContext } from "./FileField";
import TablerFile from "~icons/tabler/file";

export const FileInput = (): JSX.Element => {
  const inputRef = useRef<HTMLInputElement>(null);
  const [files, setFiles] = useState<File[]>([]);
  const [props] = useContextProps({} as InputProps, null, InputContext);
  const customProps = useContext(CustomPropsContext);

  const onPress = useCallback(() => {
    if (!inputRef.current) return;
    if (inputRef.current.value !== "") inputRef.current.value = "";
    inputRef.current.click();
  }, []);

  const onChange = useCallback((event: ChangeEvent<HTMLInputElement>) => {
    if (props.multiple) {
      setFiles((prevFiles) => {
        if (!event.target.files) {
          customProps.onChange?.([...prevFiles]);
          return [...prevFiles];
        }

        if (customProps.maxSize) {
          for (const file of event.target.files) {
            if (file.size > customProps.maxSize) {
              return [...prevFiles];
            }
          }
        }

        const newFiles = [...prevFiles, ...event.target.files];
        customProps.onChange?.(newFiles);
        return newFiles;
      });
    } else {
      if (!event.target.files) {
        setFiles([]);
        customProps.onChange?.([]);
        return;
      }

      const file = event.target.files[0];

      if (customProps.maxSize && file.size > customProps.maxSize) {
        return;
      }

      setFiles([file]);
      customProps.onChange?.([file]);
    }
  }, []);

  const onRemoveFileByIndex = (fileToRemoveIdx: number) => {
    const newFiles = files ? [...files] : [];
    newFiles.splice(fileToRemoveIdx, 1);

    updateFilesForInput(newFiles);
    setFiles(newFiles);

    customProps.onChange?.(newFiles);
  };

  // not sure if we need this method if we set validationBehavior to "aria" inside the FileField component, but it should work for now.
  const updateFilesForInput = (newFiles: File[]) => {
    if (inputRef.current) {
      const dataTransfer = new DataTransfer();
      newFiles.forEach((file) => dataTransfer.items.add(file));
      inputRef.current.files = dataTransfer.files;
    }
  };

  const allowedFileTypes = props.accept?.split(",");

  return (
    <div>
      {/* TODO: make dropzone optional */}
      <DropZone
        className={cn(
          "flex items-center justify-center border-1 border-input-border rounded-sm px-5 py-7 bg-surface min-h-25",
          inputStyles({ isInvalid: !!props["aria-invalid"] }),
        )}
        onDrop={async (e) => {
          const fileDropItems = e.items.filter(
            (file) => file.kind === "file",
          ) as FileDropItem[];

          let newFiles: File[];
          if (props.multiple) {
            const allUploadedFiles = await Promise.all(
              fileDropItems.map((fileDropItem) => fileDropItem.getFile()),
            );

            if (customProps.maxSize) {
              for (const file of allUploadedFiles) {
                if (allowedFileTypes && !allowedFileTypes.includes(file.type)) {
                  return;
                }

                if (file.size > customProps.maxSize) {
                  return;
                }
              }
            }

            newFiles = [...files, ...allUploadedFiles];
            setFiles(newFiles);
          } else {
            const file = await fileDropItems[0].getFile();

            if (allowedFileTypes && !allowedFileTypes.includes(file.type)) {
              return;
            }

            if (customProps.maxSize && file.size > customProps.maxSize) {
              return;
            }

            newFiles = [file];
            setFiles(newFiles);
          }

          customProps.resetValidation();
          updateFilesForInput(newFiles);

          customProps.onChange?.(newFiles);
        }}
      >
        <div className="flex flex-col gap-1 min-w-0">
          <Button
            className="gap-2"
            variant="secondary"
            size="sm"
            onPress={onPress}
          >
            <TablerUpload />
            Browse files
          </Button>
          <p className="text-sm text-text-secondary text-center">
            or drop files here.
          </p>
        </div>

        <Input
          style={{ display: "none" }}
          custom-hidden="true"
          type="file"
          ref={inputRef}
          onChange={onChange}
        />
      </DropZone>

      <div>
        {files && files.length > 0 && (
          <div className="mt-1 flex flex-col gap-1">
            {files.map((file, i) => (
              <div
                key={i}
                className="flex justify-between items-center gap-5 border-1 border-input-border rounded-sm p-2 bg-surface"
              >
                <div className="flex gap-1 items-center">
                  <TablerFile className="text-text-secondary" />
                  <p className="font-medium truncate text-sm">{file.name}</p>
                </div>

                <TablerX
                  onClick={() => onRemoveFileByIndex(i)}
                  className="text-sm text-gray-500 hover:text-red-500 cursor-pointer"
                />
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  );
};
