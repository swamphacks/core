import { useState } from "react";

import "./App.css";

function App() {
  return (
    <div className="container">
      <Background />

      <div className="water-strip"></div>

      <div className="fishing-line">
        <img src="./assets/string.png" className="string-img" />
        <img src="./assets/bobber.png" className="bobber-img" />
      </div>
    </div>
  );
}

function Background() {
  return (
    <div className="banner">
      <div className="layer layer-sky"></div>
      <div className="layer layer-bushes-back"></div>
      <div className="layer layer-trees"></div>
      <div className="layer layer-trees-back1"></div>
      <div className="layer layer-trees-back2"></div>
      <div className="layer layer-trees-back3"></div>
      <div className="layer layer-vegetation"></div>
      <div className="layer layer-canopy"></div>
      <div className="layer layer-canopy-back"></div>
      <div className="boat-rig">
        <img src="./assets/boat.png" className="boat" />
        <img src="./assets/rod.png" className="rod" />
      </div>
    </div>
  );
}

export default App;
