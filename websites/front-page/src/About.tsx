import "./About.css";

export default function About() {
  return (
    <div className="about-container">
      <div className="about-header">
        <h1>About Us</h1>
        <div className="about-summary-container">
          <p className="summary">
            SwampHacks is the University of Florida’s flagship hackathon,
            bringing together 400+ students each year for 36 hours of
            creativity, collaboration, and innovation. Recognized for excellence
            by UF’s Herbert Wertheim College of Engineering, SwampHacks offers
            hands-on workshops, mentorship, and community-building activities
            that help hackers grow their skills and bring their ideas to life.
            Whether you’re a first-time hacker or a seasoned coder, SwampHacks
            is the place to build, connect, and inspire.
          </p>

          <div className="stats">
            <span className="nes-badge">
              <span className="is-primary">400+ hackers</span>
            </span>
            <span className="nes-badge">
              <span className="is-primary">10+ workshops</span>
            </span>
            <span className="nes-badge">
              <span className="is-primary">100+ projects</span>
            </span>
            <span className="nes-badge">
              <span className="is-primary">$10k+ in prizes</span>
            </span>
          </div>
        </div>
      </div>
    </div>
  );
}
