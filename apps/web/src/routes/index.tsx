import {
  createFileRoute,
  redirect,
  useLocation,
  useSearch,
} from "@tanstack/react-router";
import { Login } from "@/modules/Auth/Login";
import { PageLoading } from "@/components/PageLoading";
import { z } from "zod";

export const Route = createFileRoute("/")({
  validateSearch: z.object({
    redirect: z.string().optional().catch(""),
  }),
  beforeLoad: async ({ context }) => {
    const { user } = await context.userQuery.promise;

    console.log("Loaded user in beforeLoad:", user);

    if (user) {
      console.log("User is already authenticated, redirecting to portal.");
      throw redirect({
        to: "/portal",
      });
    }
  },
  pendingMs: 50,
  // TODO: Use skeleton components to display a loading state
  pendingComponent: () => PageLoading(),
  component: Index,
});

function Index() {
  const location = useLocation();
  const search = useSearch({ from: "/" });
  const redirectTarget = search.redirect;

  console.log("current URL", location.href);
  console.log("redirect: ", redirectTarget);

  const isRedirectApplication = redirectTarget === "application";

  return (
    <div className="h-full flex items-center justify-center pb-15">
      {/* get url, check if it redirect = application, if true, show user message */}
      {isRedirectApplication && (
        <p className="mb-4 text-center text-lg font-semibold text-yellow-600">
          You must login with your Discord account first before applying!
        </p>
      )}
      <Login />
    </div>
  );
}
