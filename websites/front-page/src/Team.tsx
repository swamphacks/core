import "./Team.css";

import woodenSign from "./assets/WoodenSign.png";
import tent from "./assets/Tent.png";
import tent2 from "./assets/Tent2.png";
import fire from "./assets/Fire.png";

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
        </div>

        <div className="tent-container">
          <img className="tent" src={tent} />
          <img className="tent2" src={tent2} />
          <img className="fire" src={fire} />
        </div>
      </div>
    </div>
  );
}
