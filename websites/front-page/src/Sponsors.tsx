import type { PointerEvent } from "react";
import "./Sponsors.css";
import boat from "./assets/boat.png";
// import rod from "./assets/frod.png";
import fish from "./assets/fish.png";
// import SponsorsImg from "./assets/Sponsors.png";
// import SponsorsSvg from "./assets/Sponsors.svg";

// const sponsors = [
//   {
//     name: "Placeholder",
//     logo: UF,
//     url: "",
//     gridSize: "co-host",
//   },
//   {
//     name: "Placeholder",
//     logo: FifthThird,
//     url: "",
//     gridSize: "rect-long",
//   },
//   {
//     name: "Placeholder",
//     logo: UF,
//     url: "",
//     gridSize: "rect-long",
//   },
//   {
//     name: "Placeholder",
//     logo: UF,
//     url: "",
//     gridSize: "small",
//   },
//   {
//     name: "Placeholder",
//     logo: UF,
//     url: "",
//     gridSize: "small",
//   },
//   {
//     name: "Placeholder",
//     logo: UF,
//     url: "",
//     gridSize: "small",
//   },
// ];

export default function Sponsors() {
  const handleFishFlee = (event: PointerEvent<HTMLDivElement>) => {
    if (event.pointerType === "touch") return;

    const avoidanceRadius = 150;
    const maximumDistance = 100;
    const fishElements =
      event.currentTarget.querySelectorAll<HTMLElement>(".fish");

    fishElements.forEach((fishElement) => {
      const bounds = fishElement.getBoundingClientRect();
      const fishX = bounds.left + bounds.width / 2;
      const fishY = bounds.top + bounds.height / 2;
      const differenceX = fishX - event.clientX;
      const differenceY = fishY - event.clientY;
      const distance = Math.hypot(differenceX, differenceY);

      if (distance < avoidanceRadius && distance > 0) {
        const strength = 1 - distance / avoidanceRadius;
        const fleeX =
          (differenceX / distance) * strength * maximumDistance;
        const fleeY =
          (differenceY / distance) * strength * maximumDistance;

        const angle = Math.max(-10, Math.min(10, fleeY * 0.08));

        fishElement.style.transform =
          `translate3d(${fleeX}px, ${fleeY}px, 0) rotate(${angle}deg)`;
      } else {
        fishElement.style.transform = "translate3d(0, 0, 0) rotate(0deg)";
      }
    });
  };

  const resetFishFlee = (event: PointerEvent<HTMLDivElement>) => {
    event.currentTarget
      .querySelectorAll<HTMLElement>(".fish")
      .forEach((fishElement) => {
        fishElement.style.transform = "translate3d(0, 0, 0) rotate(0deg)";
      });
  };

  const fishSchool = [] as Array<{
    id: number;
    left: number;
    top: number;
    size: number;
    flip: boolean;
    opacity: number;
    duration: number;
    delay: number;
  }>;

  for (let index = 0; index < 15; index += 1) {
    let left = Math.random() * 90;
    let top = 20 + Math.random() * 55;
    let attempts = 0;

    while (attempts < 60) {
      const tooClose = fishSchool.some((existingFish) => {
        const horizontalDistance = Math.abs(existingFish.left - left);
        const verticalDistance = Math.abs(existingFish.top - top);
        return Math.sqrt(horizontalDistance ** 2 + verticalDistance ** 2) < 16;
      });

      if (!tooClose) {
        break;
      }

      left = Math.random() * 90;
      top = 20 + Math.random() * 55;
      attempts += 1;
    }

    fishSchool.push({
      id: index,
      left,
      top,
      size: 25 + Math.random() * 25,
      flip: Math.random() > 0.5,
      opacity: 0.6 + Math.random() * 0.35,
      duration: 16 + Math.random() * 12,
      delay: -(Math.random() * 28),
    });
  }

  return (
    <div id="sponsors" className="sponsors-container">
      <div className="sponsors-vegetation" />

      <div
        className="water"
        onPointerMove={handleFishFlee}
        onPointerLeave={resetFishFlee}
      >
        <div className="water-texture" aria-hidden="true" />

        <div
          style={{
            position: "relative",
            zIndex: 3,
            display: "flex",
            flexDirection: "column",
            alignItems: "center",
            width: "100%",
          }}
        >
          <h1 className="sponsors-header">Sponsors</h1>
          <p className="sponsors-coming-soon">Coming soon...</p>
        </div>

        <div className="fish-school" aria-hidden="true">
          {fishSchool.map((fishItem) => (
            <div
              key={fishItem.id}
              className={
                fishItem.flip
                  ? "fish-path fish-path--left"
                  : "fish-path fish-path--right"
              }
              style={{
                top: `${fishItem.top}%`,
                animationDuration: `${fishItem.duration}s`,
                animationDelay: `${fishItem.delay}s`,
              }}
            >
              <img
                className="fish"
                src={fish}
                alt=""
                style={{
                  width: `${fishItem.size}px`,
                  opacity: fishItem.opacity,
                }}
              />
            </div>
          ))}
        </div>

        <div className="boat-outer-container">
          <div className="boat-container">
            <img className="boat" src={boat} alt="SwampHacks ship" />
            <img
              className="boat boat-reflection"
              src={boat}
              alt=""
              aria-hidden="true"
            />
          </div>
        </div>
      </div>
    </div>
  );
}
