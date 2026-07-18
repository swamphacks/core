import "./Team.css";

import woodenSign from "./assets/WoodenSign.png";
import tent from "./assets/Tent.png";
import tent2 from "./assets/Tent2.png";
import fire from "./assets/Fire.png";
import pfp from "./assets/shxii-icon-transparent.png"
import glow from "./assets/glowing.png"

export default function Team() {
  return (
    <div className="team-container">
      <h1 className="team-header">Meet the XII Team!</h1>

      {/* <div className="background2">
        <div className="foreground2"></div>
      </div> */}

      <div className="team-area">
        <div className="sign-container">
          <img className="sign" src={woodenSign} />
          <img className="pfp" src={pfp} />
        </div>

        <div className="tent-container">
          <img className="tent" src={tent} />
          <img className="tent2" src={tent2} />
          <img className="fire" src={fire} />
        </div>
      </div>

      <div className="fly-container">
        <img
          className="fly"
          src={pfp}
          style={{
            top: "205px",
            left: "20.5%",
          }}
        />
        <img
          className="glow"
          src={glow}
          style={{
            top: "200px",
            left: "20%",
          }}
        />
      </div>
    </div>
  );
}
