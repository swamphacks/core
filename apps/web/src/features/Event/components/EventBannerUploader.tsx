import { Button } from "@/components/ui/Button";
import { useState } from "react";
import { Button as RAC_Button, FileTrigger } from "react-aria-components";
import { useEventBannerActions } from "../hooks/useEventBannerActions";
import { toast } from "react-toastify";
import imageCompression, { type Options } from "browser-image-compression";

interface Props {
  bannerUrl: string | null;
  eventId: string;
}

const EventBannerUploader = ({ bannerUrl, eventId }: Props) => {
  // State to hold the selected file
  const [file, setFile] = useState<File | null | undefined>(undefined); // File and null can be uploaded, undefined means no change
  const [currentBanner, setCurrentBanner] = useState<string | null>(bannerUrl);
  const [uploadProgress, setUploadProgress] = useState(0); // 0-100
  const [isUploading, setIsUploading] = useState(false);

  const { upload, remove } = useEventBannerActions(eventId);

  const handleFileSelection = async (e: FileList | null) => {
    const files = e ? Array.from(e) : [];

    if (!files[0]) {
      setFile(null);
      return;
    }

    const selectedFile = files[0];

    const options = {
      maxSizeMB: 0.5,
      maxWidthOrHeight: 1920,
      useWebWorker: true,
    } as Options;

    // Compress the image file
    const compressedFileBlob = await imageCompression(selectedFile, options);

    // Create a new File object from the compressed Blob to preserve the file name and type
    // This is needed so the server can correctly identify the file type
    const compressedFile = new File([compressedFileBlob], selectedFile.name, {
      type: compressedFileBlob.type,
    });

    setFile(compressedFile);
  };

  const handleUpload = async () => {
    if (file) {
      setIsUploading(true);
      setUploadProgress(0);

      // Fake progress bar until upload completes
      let progress = 0;
      const interval = setInterval(() => {
        progress += Math.random() * 3 + 1; // Random increment between 1-4%
        if (progress >= 80) progress = 80; // Cap at 80% until upload finishes
        setUploadProgress(progress);
      }, 10); // Every 10ms to fake a smooth progress

      await upload.mutateAsync(file, {
        onSuccess: (data) => {
          clearInterval(interval);
          setUploadProgress(100);
          setTimeout(() => {
            setIsUploading(false);
            setUploadProgress(0);
            setCurrentBanner(data.banner_url);
            setFile(undefined);
          }, 500);
          toast.success("Banner uploaded successfully!", {
            position: "bottom-right",
          });
        },
        onError: () => {
          clearInterval(interval);
          setIsUploading(false);
          setUploadProgress(0);
          toast.error("Failed to upload banner. Please try again.", {
            position: "bottom-right",
          });
        },
      });
    } else {
      await remove.mutateAsync(undefined, {
        onSuccess: () => {
          setCurrentBanner(null);
          toast.success("Banner removed successfully!", {
            position: "bottom-right",
          });
        },
        onError: () => {
          toast.error("Failed to remove banner. Please try again.", {
            position: "bottom-right",
          });
        },
      });
    }
  };

  return (
    <div className="max-w-lg flex flex-col gap-2">
      <p>Banner</p>
      {file === null || (file === undefined && !currentBanner) ? (
        <FileTrigger
          acceptedFileTypes={["image/*"]}
          onSelect={handleFileSelection}
        >
          <RAC_Button className="w-full h-58 flex items-center justify-center bg-neutral-100 dark:bg-neutral-800 border border-dashed border-neutral-300 dark:border-neutral-700 rounded-sm cursor-pointer select-none hover:dark:bg-neutral-700 hover:bg-neutral-200 transition-colors">
            <span className="text-neutral-500">No Banner Uploaded</span>
          </RAC_Button>
        </FileTrigger>
      ) : (
        <>
          <FileTrigger
            acceptedFileTypes={["image/png", "image/jpeg", "image/jpg"]}
            allowsMultiple={false}
            onSelect={handleFileSelection}
          >
            <RAC_Button className="relative w-full h-58 group p-0 overflow-hidden rounded-sm cursor-pointer">
              <img
                src={file ? URL.createObjectURL(file) : currentBanner || ""}
                alt="Event Banner"
                className="w-full h-58 object-cover rounded-sm"
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

      {/* Progress bar for uploads */}
      {isUploading && (
        <div className="absolute bottom-0 left-0 w-full h-1 bg-neutral-200 dark:bg-neutral-600 rounded-t overflow-hidden">
          <div
            className="h-1 bg-cyan-700 transition-[width] duration-300 ease-out"
            style={{ width: `${uploadProgress}%` }}
          />
        </div>
      )}
    </div>
  );
};

export default EventBannerUploader;
