import { createFileRoute } from "@tanstack/react-router";
import { useEffect, useState } from "react";
import { OnboardingModal } from "@/modules/Onboarding/OnboardingModal";
import Cookies from "js-cookie";
import { useTheme } from "@/components/ThemeProvider";
import {
  Timeline,
  TimelineContent,
  TimelineDate,
  TimelineHeader,
  TimelineIndicator,
  TimelineItem,
  TimelineSeparator,
  TimelineTitle,
} from "@/components/reui/timeline";

const items = [
  {
    date: "May 2024",
    sortDate: "2024-05-01",
    title: "Project Initialized",
    content:
      "Successfully set up the project repository and initial architecture.",
    iconDark: "../../../public/assets/SwampHacks_Logo_Dark.png",
    iconLight: "../../../public/assets/SwampHacks_Logo_Light.png",
  },
  {
    date: "April 2024",
    sortDate: "2024-04-01",
    title: "Beta Release",
    content: "Launched the beta version for early testers and feedback.",
    iconDark: "../../../public/assets/SwampHacks_Logo_Dark.png",
    iconLight: "../../../public/assets/SwampHacks_Logo_Light.png",
  },
  {
    date: "June 2024",
    sortDate: "2024-06-01",
    title: "Official Launch",
    content: "The platform is now live for all users worldwide.",
    iconDark: "../../../public/assets/SwampHacks_Logo_Dark.png",
    iconLight: "../../../public/assets/SwampHacks_Logo_Light.png",
  },
  {
    date: "May 2025",
    sortDate: "2025-05-01",
    title: "ex3",
    content: "bloh",
    iconDark: "../../../public/assets/SwampHacks_Logo_Dark.png",
    iconLight: "../../../public/assets/SwampHacks_Logo_Light.png",
  },
  {
    date: "May 2026",
    sortDate: "2026-05-01",
    title: "ex4",
    content: "blah",
    iconDark: "../../../public/assets/SwampHacks_Logo_Dark.png",
    iconLight: "../../../public/assets/SwampHacks_Logo_Light.png",
  },
  {
    date: "May 2027",
    sortDate: "2027-05-01",
    title: "ex5",
    content: "bleh",
    iconDark: "../../../public/assets/SwampHacks_Logo_Dark.png",
    iconLight: "../../../public/assets/SwampHacks_Logo_Light.png",
  },
];

const sortedItems = [...items].sort(
  (a, b) => new Date(a.sortDate).getTime() - new Date(b.sortDate).getTime(),
);

function getCurrentStep(items: typeof sortedItems) {
  const now = Date.now();
  let activeIndex = 0;
  for (let i = 0; i < items.length; i++) {
    if (new Date(items[i].sortDate).getTime() <= now) {
      activeIndex = i;
    } else {
      break;
    }
  }
  return activeIndex + 1;
}

export function Pattern() {
  const { theme } = useTheme();
  const isDark = theme === "dark";
  const [visibleCount, setVisibleCount] = useState(0);
  const [activeStep, setActiveStep] = useState(0);

  const targetStep = getCurrentStep(sortedItems);
  const buildComplete = visibleCount === sortedItems.length;

  useEffect(() => {
    if (buildComplete) return;

    const timeout = setTimeout(() => {
      setVisibleCount((c) => c + 1);
    }, 400);

    return () => clearTimeout(timeout);
  }, [visibleCount, buildComplete]);

  useEffect(() => {
    if (!buildComplete) return;
    if (activeStep >= targetStep) return;

    const delay = activeStep === 0 ? 300 : 350;
    const timeout = setTimeout(() => {
      setActiveStep((s) => s + 1);
    }, delay);

    return () => clearTimeout(timeout);
  }, [buildComplete, activeStep, targetStep]);

  const visibleItems = sortedItems.slice(0, visibleCount);

  return (
    <Timeline value={activeStep} className="w-full max-w-md">
      {visibleItems.map((item, index) => {
        const isCompleted = index + 1 <= activeStep;

        const activeIcon = isDark ? item.iconDark : item.iconLight;
        const inactiveIcon = isDark ? item.iconLight : item.iconDark;
        const iconSrc = isCompleted ? activeIcon : inactiveIcon;

        return (
          <TimelineItem
            key={item.sortDate + item.title}
            step={index + 1}
            className="animate-in fade-in slide-in-from-bottom-2 duration-500"
          >
            <TimelineHeader>
              <TimelineDate>{item.date}</TimelineDate>
              <TimelineTitle>{item.title}</TimelineTitle>
            </TimelineHeader>
            <TimelineIndicator className="flex items-center justify-center bg-surface overflow-hidden p-0">
              <img
                src={iconSrc}
                alt=""
                className="absolute inset-0 size-full object-cover rounded-full transition-opacity duration-300"
              />
            </TimelineIndicator>
            <TimelineSeparator />
            <TimelineContent>{item.content}</TimelineContent>
          </TimelineItem>
        );
      })}
    </Timeline>
  );
}

export const Route = createFileRoute("/_protected/information")({
  beforeLoad: (context) => {
    const { user } = context.context;
    const hasSkippedCookie = Cookies.get("welcome-modal-skipped") === "true";
    const showOnboardingModal = !hasSkippedCookie && !!user && !user.onboarded;
    return {
      showOnboardingModal: showOnboardingModal,
    };
  },
  component: RouteComponent,
});

function RouteComponent() {
  const { showOnboardingModal } = Route.useRouteContext();
  const [isModalOpen, setIsModalOpen] = useState(showOnboardingModal);

  return (
    <div className="min-h-screen bg-surface text-text-main overflow-x-hidden -m-6 pt-10 pb-10">
      <div className="max-w-2xl mx-auto">
        <span className="block text-4xl font-bold text-text-main mb-1">
          Welcome!
        </span>
        <span className="block text-2xl font-bold text-text-main mb-1">
          SwampHacks XII Oct 10-13
        </span>
        <span className="block w-full pb-2 text-xs italic text-text-secondary">
          Current schedule which is subject to some change
        </span>
        <div className="p-5 border-border border-1 shadow-xl shadow-ring rounded-lg">
          <Pattern />
        </div>
      </div>
      <OnboardingModal isOpen={isModalOpen} onOpenChange={setIsModalOpen} />
    </div>
  );
}
