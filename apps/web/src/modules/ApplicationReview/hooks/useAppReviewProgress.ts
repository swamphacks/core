import { useState, useEffect } from "react";
import type { AssignedApplications } from "./useAssignedApplications";

/**
 * Custom hook to manage application review progress.
 * Supports navigating through assigned applications and
 * a virtual "Review Completed" step at the end.
 *
 * @param apps - Array of assigned applications
 */
export const useAppReviewProgress = (apps: AssignedApplications) => {
  /**
   * Determines the initial index when the hook is first used.
   * - If apps are empty → -1 (loading / no apps)
   * - If all apps completed → virtual finish step (apps.length)
   * - If any in-progress app → index of first in-progress
   * - Otherwise → fallback to first app (index 0)
   */
  const getInitialIndex = () => {
    if (apps.length === 0) return -1;

    const allCompleted = apps.every((app) => app.status === "completed");
    if (allCompleted) return apps.length;

    const inProgressIndex = apps.findIndex(
      (app) => app.status === "in_progress",
    );
    return inProgressIndex >= 0 ? inProgressIndex : 0;
  };

  const [currentIndex, setCurrentIndex] = useState<number>(getInitialIndex());

  /**
   * Effect to keep currentIndex valid when apps change
   * Handles edge cases such as:
   * - Apps being loaded asynchronously
   * - Apps being updated or removed
   * - Initial load
   */
  useEffect(() => {
    if (apps.length === 0) {
      setCurrentIndex(-1);
      return;
    }

    // If currentIndex is invalid (negative or past virtual finish), reset it
    if (currentIndex < 0 || currentIndex > apps.length) {
      const newIndex = apps.findIndex((app) => app.status === "in_progress");
      setCurrentIndex(newIndex >= 0 ? newIndex : apps.length);
    }
  }, [apps, currentIndex]);

  /**
   * The application corresponding to the current index.
   * Null when at the virtual finish step.
   */
  const currentAssignedApplication =
    currentIndex < apps.length ? apps[currentIndex] : null;

  /**
   * Move to the next step.
   * Caps at virtual finish step (apps.length).
   */
  const goNext = () => {
    setCurrentIndex((prev) => Math.min(prev + 1, apps.length));
  };

  /**
   * Move to the previous step.
   * Minimum is 0.
   */
  const goPrevious = () => {
    setCurrentIndex((prev) => Math.max(prev - 1, 0));
  };

  /**
   * Jump to a specific index.
   * Allows the virtual finish step (apps.length).
   */
  const goToIndex = (index: number) => {
    if (index < 0 || index > apps.length) return;
    setCurrentIndex(index);
  };

  /**
   * Whether the user is currently on the virtual finish step.
   */
  const finished = currentIndex === apps.length;

  return {
    currentAssignedApplication,
    currentIndex,
    totalApplications: apps.length,
    finished,
    goNext,
    goPrevious,
    goToIndex,
  };
};
