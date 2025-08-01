import { type FormProps, Form as RACForm } from "react-aria-components";
import { cn } from "@/utils/cn";
import { Suspense } from "react";
import { Spinner } from "@/components/ui/Spinner";

export function Form(props: FormProps) {
  return (
    // We use suspense here since form fields are lazy loaded
    <Suspense
      fallback={
        <div className="flex w-full justify-center pt-10 gap-2 text-text-secondary">
          <Spinner />
          <p>Loading form...</p>
        </div>
      }
    >
      <RACForm
        {...props}
        className={cn("flex flex-col gap-4", props.className)}
      />
    </Suspense>
  );
}
