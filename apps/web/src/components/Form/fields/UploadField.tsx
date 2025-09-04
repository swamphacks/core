import { FileField, type FileFieldProps } from "@/components/ui/FileField";
import { useFieldContext } from "@/components/Form/formContext";

export default function UploadField(props: FileFieldProps) {
  const field = useFieldContext();

  return (
    <FileField
      onChange={(files) =>
        field.handleChange(files.length === 0 ? undefined : files)
      }
      {...props}
    />
  );
}
