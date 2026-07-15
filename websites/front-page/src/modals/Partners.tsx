import UF from "../assets/sponsors/uf.png";

interface PartnersProps {
  closeModal: () => void;
}

const sponsors = [
  {
    name: "Placeholder",
    logo: UF,
    url: "",
    gridSize: "co-host",
  },
  {
    name: "Placeholder",
    logo: UF,
    url: "",
    gridSize: "rect-long",
  },
  {
    name: "Placeholder",
    logo: UF,
    url: "",
    gridSize: "rect-long",
  },
  {
    name: "Placeholder",
    logo: UF,
    url: "",
    gridSize: "small",
  },
  {
    name: "Placeholder",
    logo: UF,
    url: "",
    gridSize: "small",
  },
  {
    name: "Placeholder",
    logo: UF,
    url: "",
    gridSize: "small",
  },
];

export default function Partners({ closeModal }: PartnersProps) {
  return (
    <div className="sponsor-modal-container">
      <button onClick={closeModal} className="modal-close-btn nes-btn is-error">
        X
      </button>
      <div>
        <p className="tab-title">Our Community Partners</p>
        <div className="tab-container sponsors-container">
          <div className="sponsors-grid">
            {sponsors.map((sponsor) => (
              <div
                className={`group sponsor-card sponsor-${sponsor.gridSize} gentle-float`}
              >
                <a
                  // href={sponsor.url}
                  target="_blank"
                  rel="noopener noreferrer"
                  // aria-label={sponsor.name}
                  className="relative block h-full"
                >
                  <div className="sponsor-logo-card">
                    <img
                      src={sponsor.logo}
                      className={`sponsor-logo logo-medium`}
                    />
                  </div>
                  {/* <div className="sponsor-tooltip">Visit {sponsor.name}</div> */}
                </a>
              </div>
            ))}
          </div>
          {/* <div className="sponsor">
            <img src={Vobile} />
          </div> */}
        </div>
      </div>
    </div>
  );
}
