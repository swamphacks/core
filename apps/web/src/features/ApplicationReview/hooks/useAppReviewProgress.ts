import { useState, useEffect } from "react";
import type { AssignedApplications } from "./useAssignedApplications";

export const useAppReviewProgress = (apps: AssignedApplications) => {
  const [currentIndex, setCurrentIndex] = useState(() =>
    apps.findIndex((app) => app.status === "in_progress"),
  );
  const [finished, setFinished] = useState(() =>
    apps.every((app) => app.status === "completed"),
  );

  // Keep currentIndex valid when apps update
  useEffect(() => {
    if (!apps.length) {
      setCurrentIndex(-1);
      return;
    }

    // If current index is invalid or current app removed, reset to first "in_progress" or first app
    if (currentIndex < 0 || currentIndex >= apps.length) {
      const newIndex = apps.findIndex((app) => app.status === "in_progress");
      setCurrentIndex(newIndex >= 0 ? newIndex : 0);
    }

    setFinished(apps.every((app) => app.status === "completed"));
  }, [apps, currentIndex]);

  const currentAssignedApplication = apps[currentIndex] || null;

  const goNext = () => {
    setCurrentIndex((prev) => Math.min(prev + 1, apps.length - 1));

    if (currentIndex + 1 >= apps.length) {
      setFinished(true);
    }
  };
  const goPrevious = () => setCurrentIndex((prev) => Math.max(prev - 1, 0));

  return {
    currentAssignedApplication,
    goNext,
    goPrevious,
    currentIndex,
    finished,
    totalApplications: apps.length,
  };
};
