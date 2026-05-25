import "./App.css";
import swamphacksIcon from "./assets/shxii-icon-transparent.png";
import gatorGif from "./assets/gator.gif";

function App() {
  return (
    <div className="container">
      <main>
        <header className="hero">
          <div>
            <img src={swamphacksIcon} className="icon-image" />
          </div>
          <h1 className="title">
            SwampHacks XII
            <span className="badge">Apply now!</span>
          </h1>
          <h1 className="subtitle">Oct 16 - 18, 2026 | Reitz Union</h1>
        </header>

        <div className="menu-container">
          <button className="nes-btn is-primary">Register</button>
          <button className="nes-btn is-primary">About</button>
          <div className="community-buttons">
            <button className="nes-btn is-primary">Sponsors</button>
            <button className="nes-btn is-primary">Partners</button>
          </div>
        </div>

        <section className="icon-list">
          <a href="https://www.instagram.com/ufswamphacks/" target="_blank">
            <i className="nes-icon instagram"></i>
          </a>

          <a href="https://discord.gg/YBHrXPJ8mR" target="_blank">
            <i className="nes-icon discord"></i>
          </a>

          <a href="https://github.com/swamphacks/core" target="_blank">
            <i className="nes-icon github"></i>
          </a>
        </section>

        <img
          src={gatorGif}
          alt="Pixel-art alligator walking in front of the sign"
          className="gator"
        />
      </main>

      <footer>
        <p>Made with ❤️ by the SwampHacks Team © 2026</p>
      </footer>
    </div>
  );
}

export default App;
