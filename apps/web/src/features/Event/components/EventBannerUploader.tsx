import { Button } from "@/components/ui/Button";
import { useState } from "react";
import { Button as RAC_Button, FileTrigger } from "react-aria-components";
import { useUploadEventBanner } from "../hooks/useUploadEventBanner";
import { toast } from "react-toastify";

interface Props {
  bannerUrl: string | null;
  eventId: string;
}

const EventBannerUploader = ({ bannerUrl, eventId }: Props) => {
  // State to hold the selected file
  const [file, setFile] = useState<File | null | undefined>(undefined); // File and null can be uploaded, undefined means no change
  const { mutateAsync } = useUploadEventBanner(eventId);

  const handleFileSelection = (e: FileList | null) => {
    let files = e ? Array.from(e) : [];
    setFile(files[0] || null);
  };

  const handleUpload = async () => {
    if (file) {
      console.log("Uploading file:", file);
      await mutateAsync(file, {
        onSuccess: (data) => {
          bannerUrl = data.banner_url;
          setFile(undefined); // Reset file state after successful upload
          toast.success("Banner uploaded successfully!", {
            position: "bottom-right",
          });
        },
        onError: () => {
          toast.error("Failed to upload banner. Please try again.", {
            position: "bottom-right",
          });
        },
      });
    } else {
      toast.info("No file selected for upload.", {
        position: "bottom-right",
      });
    }
  };

  return (
    <div className="max-w-lg flex flex-col gap-2">
      <p>Banner</p>
      {bannerUrl || file ? (
        <>
          <FileTrigger
            acceptedFileTypes={["image/png", "image/jpeg", "image/jpg"]}
            allowsMultiple={false}
            onSelect={handleFileSelection}
          >
            <RAC_Button className="relative w-full h-58 group p-0 overflow-hidden rounded-sm cursor-pointer">
              <img
                src={bannerUrl || (file && URL.createObjectURL(file)) || ""}
                alt="Event Banner"
                className="object-cover w-full h-full rounded-sm border border-neutral-300 dark:border-neutral-700"
              />
              {/* Overlay */}
              <div className="absolute inset-0 bg-black/60 opacity-0 group-hover:opacity-100 transition-opacity flex flex-row justify-center items-center">
                <span className="text-white">Change Banner</span>
              </div>
            </RAC_Button>
          </FileTrigger>
          <Button
            variant="danger"
            size="sm"
            className="rounded-sm"
            onPress={() => setFile(null)}
          >
            Remove Banner
          </Button>
        </>
      ) : (
        <FileTrigger
          acceptedFileTypes={["image/*"]}
          onSelect={handleFileSelection}
        >
          <RAC_Button
            className="w-full h-58 flex items-center justify-center bg-neutral-100 dark:bg-neutral-800 border border-dashed border-neutral-300 dark:border-neutral-700 rounded-sm cursor-pointer select-none hover:dark:bg-neutral-700 hover:bg-neutral-200 transition-colors"
          >
            <span className="text-neutral-500">No Banner Uploaded</span>
          </RAC_Button>
        </FileTrigger>
      )}
      {file !== undefined && (
        <Button
          variant="primary"
          size="sm"
          className="rounded-sm"
          onPress={handleUpload}
        >
          Save Banner Changes
        </Button>
      )}
      <p className="text-sm text-neutral-500">
        Recommended dimensions: 1200x300px. Max file size: 5MB. Supported
        formats: JPG, JPEG, PNG.
      </p>
    </div>
  );
};

export default EventBannerUploader;
