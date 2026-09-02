import { useState, useEffect } from "react";
import "./Countdown.css";

interface TimeLeft {
  days: number;
  hours: number;
  minutes: number;
  seconds: number;
  done: boolean;
}

interface CountdownCornerProps {
  targetDate?: Date;
  label?: string;
}

function getTimeLeft(targetDate: Date): TimeLeft {
  const diff = +targetDate - +new Date();
  if (diff <= 0) {
    return { days: 0, hours: 0, minutes: 0, seconds: 0, done: true };
  }
  return {
    days: Math.floor(diff / (1000 * 60 * 60 * 24)),
    hours: Math.floor((diff / (1000 * 60 * 60)) % 24),
    minutes: Math.floor((diff / 1000 / 60) % 60),
    seconds: Math.floor((diff / 1000) % 60),
    done: false,
  };
}

export default function CountdownCorner({
  targetDate = new Date(Date.now() + 1000 * 60 * 60 * 24 * 3),
  label = "Swamphacks Begins In:",
}: CountdownCornerProps) {
  const [timeLeft, setTimeLeft] = useState<TimeLeft>(() =>
    getTimeLeft(targetDate),
  );

  useEffect(() => {
    const interval = setInterval(() => {
      setTimeLeft(getTimeLeft(targetDate));
    }, 1000);
    return () => clearInterval(interval);
  }, [targetDate]);

  if (timeLeft.done) return null;

  const pad = (n: number) => String(n).padStart(2, "0");

  return (
    <div className="countdown-corner">
      <span className="countdown-butterfly" aria-hidden="true">
        <span className="countdown-butterfly__body" />
      </span>
      <span className="countdown-fireflies" aria-hidden="true">
        {Array.from({ length: 5 }).map((_, index) => (
          <span className="countdown-firefly" key={index} />
        ))}
      </span>
      <span className="countdown-corner__label">{label}</span>
      <div className="countdown-corner__digits">
        <span className="countdown-corner__unit">
          <strong>{pad(timeLeft.days)}</strong>
          <small>days</small>
        </span>
        <span className="countdown-corner__separator">:</span>
        <span className="countdown-corner__unit">
          <strong>{pad(timeLeft.hours)}</strong>
          <small>hrs</small>
        </span>
        <span className="countdown-corner__separator">:</span>
        <span className="countdown-corner__unit">
          <strong>{pad(timeLeft.minutes)}</strong>
          <small>min</small>
        </span>
        <span className="countdown-corner__separator">:</span>
        <span className="countdown-corner__unit">
          <strong>{pad(timeLeft.seconds)}</strong>
          <small>sec</small>
        </span>
      </div>
    </div>
  );
}
