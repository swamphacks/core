import { type FormProps, Form as RACForm } from "react-aria-components";
import { cn } from "@/utils/cn";
import { Suspense } from "react";

export function Form(props: FormProps) {
  return (
    // We use suspense here since form fields are lazy loaded
    <Suspense fallback={<p>Loading form...</p>}>
      <RACForm
        {...props}
        className={cn("flex flex-col gap-4", props.className)}
      />
    </Suspense>
  );
}
