import { useFormContext } from "@/components/Form/formContext";
import { Button } from "@/components/ui/Button";

export function SubmitButton({ label = "Submit" }: { label?: string }) {
  const form = useFormContext();
  return (
    <form.Subscribe
      selector={(state) => [state.canSubmit, state.isSubmitting]}
      children={([canSubmit, isSubmitting]) => {
        return (
          <div className="flex justify-start">
            <Button type="submit" isDisabled={!canSubmit || isSubmitting}>
              {isSubmitting ? "Submitting..." : label}
            </Button>
          </div>
        );
      }}
    />
  );
}
