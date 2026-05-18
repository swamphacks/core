import backgroundImage from "./assets/background.png";
import boardImage from "./assets/board.png";
import gatorGif from "./assets/gator.gif";
import recapLabel from "./assets/sh xi recap.png";
import signupLabel from "./assets/sign up for updates.png";
import sponsorLabel from "./assets/sponsor sh xii.png";
import comingSoonBoard from "./assets/redboard.png";
import "./ComingSoonPage.css";

const ctaLinks = [
  {
    href: "#",
    label: "Sign up for updates",
    image: signupLabel,
    className: "coming-soon__button coming-soon__button--top",
  },
  {
    href: "#",
    label: "SwampHacks XI recap",
    image: recapLabel,
    className: "coming-soon__button coming-soon__button--middle",
  },
  {
    href: "#",
    label: "Sponsor SwampHacks XII",
    image: sponsorLabel,
    className: "coming-soon__button coming-soon__button--bottom",
  },
] as const;

export function ComingSoonPage() {
  return (
    <main
      className="coming-soon"
      style={{ backgroundImage: `url(${backgroundImage})` }}
    >
      <section className="coming-soon__content">
        <header className="coming-soon__hero">
          <h1 className="coming-soon__title">SwampHacks XII</h1>
        </header>

        <div className="coming-soon__board-wrap">
          <img
            src={boardImage}
            alt=""
            className="coming-soon__board-image"
            aria-hidden="true"
          />

          <div className="coming-soon__actions">
            {ctaLinks.map((link) => (
              <a
                key={link.label}
                href={link.href}
                className={link.className}
                aria-label={link.label}
              >
                <img src={link.image} alt="" aria-hidden="true" />
              </a>
            ))}
          </div>

          <div className="coming-soon__badge" aria-label="Coming soon">
            <img src={comingSoonBoard} alt="" aria-hidden="true" />
            <span>Coming Soon!</span>
          </div>
        </div>

        <img
          src={gatorGif}
          alt="Pixel-art alligator walking in front of the sign"
          className="coming-soon__gator"
        />
      </section>
    </main>
  );
}
