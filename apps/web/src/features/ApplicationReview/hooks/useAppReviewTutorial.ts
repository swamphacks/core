import { useEffect, useState } from "react";

/**
 * This hook manages the state of the application review tutorial.
 *
 * @returns An object containing the tutorial completion status, loading state, and a function to complete the tutorial
 */
export const useAppReviewTutorial = () => {
  const [isCompleted, setIsCompleted] = useState(false);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    setIsLoading(true);
    const completed = localStorage.getItem("appReviewTutorialCompleted");

    if (!completed) {
      setIsCompleted(false);
      setIsLoading(false);
      return;
    }

    if (completed === "true") {
      setIsCompleted(true);
    }
    setIsLoading(false);
  }, []);

  const completeTutorial = () => {
    localStorage.setItem("appReviewTutorialCompleted", "true");
    setIsCompleted(true);
  };

  return { isCompleted, isLoading, completeTutorial };
};
