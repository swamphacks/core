import { Button } from "@/components/ui/Button";
import { Group, Heading, Text } from "react-aria-components";
import TablerHome from "~icons/tabler/home";
import TablerChevronLeft from "~icons/tabler/chevron-left";
import { useRouter } from "@tanstack/react-router";

export default function NotFoundPage() {
  const router = useRouter();

  const goBackFn = router.history.canGoBack()
    ? router.history.back
    : () => router.navigate({ to: "/portal" });

  return (
    <main className="relative w-screen h-screen flex flex-col justify-center items-center px-4 sm:px-6 text-center max-w-3xl mx-auto gap-8">
      {/* Error Code */}
      <Heading className="text-8xl md:text-9xl font-medium text-neutral-400">
        404
      </Heading>

      {/* Message Block */}
      <section className="flex flex-col gap-3">
        <Text className="text-3xl">Page Not Found</Text>
        <Text className="text-lg">
          The page you&apos;re looking for doesn&apos;t exist or has been moved.
        </Text>
      </section>

      {/* Action Buttons */}
      <Group className="w-full max-w-md flex flex-col sm:flex-row gap-4 sm:gap-6 mt-2">
        <Button
          variant="primary"
          size="lg"
          className="w-full flex items-center justify-center gap-2 sm:gap-4"
          onClick={() => router.navigate({ to: "/portal" })}
        >
          <TablerHome className="w-5 h-5" />
          <span>Dashboard</span>
        </Button>

        <Button
          variant="secondary"
          size="lg"
          className="w-full flex items-center justify-center gap-2 sm:gap-4"
          onClick={() => goBackFn()}
        >
          <TablerChevronLeft className="w-5 h-5" />
          <span>Go Back</span>
        </Button>
      </Group>
    </main>
  );
}
