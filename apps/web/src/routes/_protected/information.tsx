import { createFileRoute, HeadContent } from "@tanstack/react-router";
import { useEffect, useRef, useState } from "react";
import { OnboardingModal } from "@/modules/Onboarding/OnboardingModal";
import Cookies from "js-cookie";
import { Timeline } from "@/components/ui/Timeline";
import TablerCalendarDue from "~icons/tabler/calendar-due";

declare global {
  interface Window {
    L: any;
  }
}

export const Route = createFileRoute("/_protected/information")({
  beforeLoad: (context) => {
    const { user } = context.context;
    const hasSkippedCookie = Cookies.get("welcome-modal-skipped") === "true";
    const showOnboardingModal = !hasSkippedCookie && !!user && !user.onboarded;
    return {
      showOnboardingModal: showOnboardingModal,
    };
  },
  component: RouteComponent,
});

function RouteComponent() {
  const { showOnboardingModal, user } = Route.useRouteContext();
  const [isModalOpen, setIsModalOpen] = useState(showOnboardingModal);
  const mapRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    // Original Map: https://campusmap.ufl.edu/#/

    if (!mapRef.current) return;
    if (mapRef.current.childElementCount > 0) return;

    const locationCoords = [
      [
        [29.6464695097802, -82.3469994933412],
        [29.6464697910354, -82.3470073533353],
        [29.6464827858778, -82.347007441156],
        [29.6464828520678, -82.3470232131956],
        [29.6466051573548, -82.3470217641063],
        [29.6466095388463, -82.3475589979101],
        [29.6465887165843, -82.3475661984485],
        [29.6465931066825, -82.3479455049661],
        [29.6466361172264, -82.3479492631069],
        [29.6466363875235, -82.3480118946612],
        [29.6466088900938, -82.3480120488359],
        [29.6466111851002, -82.3485777540576],
        [29.6464377164505, -82.348584451896],
        [29.6464378506758, -82.3486048005807],
        [29.6463720231527, -82.3485755311229],
        [29.6463701933914, -82.348555998975],
        [29.6462788836149, -82.3485565141597],
        [29.6462617537789, -82.3485566105048],
        [29.646261674501, -82.3482721760246],
        [29.6462616668505, -82.3482442932653],
        [29.6462616604824, -82.348221140442],
        [29.6462616027418, -82.3480166355713],
        [29.6460214567163, -82.3480179927932],
        [29.6460220961309, -82.3479534568028],
        [29.6459797386776, -82.3479551045416],
        [29.6459621998335, -82.3479557872564],
        [29.645937189167, -82.3479567605613],
        [29.6458011796817, -82.3479620528772],
        [29.6457981248406, -82.3476869073321],
        [29.6457970309735, -82.3475641187843],
        [29.6459416370515, -82.347562778174],
        [29.6459414479843, -82.3475187076878],
        [29.6459254740265, -82.3475211571131],
        [29.6459254337546, -82.3475118220426],
        [29.6459253009583, -82.3474812852922],
        [29.6459402100679, -82.3474820021915],
        [29.6459390600046, -82.3474835031581],
        [29.6459563244594, -82.3474834718947],
        [29.6459691379641, -82.3474833973893],
        [29.6459777552231, -82.3474833495427],
        [29.6459780953612, -82.3474833473554],
        [29.6459946297249, -82.3474832544006],
        [29.6459986275141, -82.3474832309984],
        [29.6460055981851, -82.3474831918861],
        [29.6460258001501, -82.3474825868422],
        [29.6460416902446, -82.3474821101212],
        [29.6460532423199, -82.3474796292121],
        [29.6460754788018, -82.3474748534877],
        [29.6460829775111, -82.3474715429172],
        [29.6460920814954, -82.347467524153],
        [29.64610705473, -82.3474609138616],
        [29.646108075939, -82.3474604631063],
        [29.6461366333515, -82.3474404900776],
        [29.6461640926009, -82.3474008751333],
        [29.646185653392, -82.3473536569378],
        [29.646195273322, -82.3473171876726],
        [29.6461990921729, -82.3472756820616],
        [29.6462001480031, -82.3472642015535],
        [29.6461999019483, -82.3472071095525],
        [29.6461996276144, -82.3471438740673],
        [29.6461992599857, -82.3470589038473],
        [29.6462167255663, -82.3470587028871],
        [29.6462443219701, -82.3470011545663],
        [29.6464695097802, -82.3469994933412],
      ],
    ];

    const polygon = window.L.polygon(locationCoords, {
      color: "#1a6bed",
      fillOpacity: 0,
      pane: "polygonPane",
    });

    const map = window.L.map(mapRef.current).setView([29.64631, -82.34788], 17);

    map.createPane("basePane");
    map.getPane("basePane").style.zIndex = "200";

    map.createPane("campusPane");
    map.getPane("campusPane").style.zIndex = "400";

    map.createPane("polygonPane");
    map.getPane("polygonPane").style.zIndex = "1000";

    map.createPane("markerPane");
    map.getPane("markerPane").style.zIndex = "1100";

    const base = window.L.esri.Vector.vectorTileLayer(
      "291da5eab3a0412593b66d384379f89f",
      {
        pane: "basePane",
        maxZoom: 20,
      },
    );

    const ufcampus = window.L.esri.Vector.vectorTileLayer(
      "https://tiles.arcgis.com/tiles/IiuFUnlkob76Az9k/arcgis/rest/services/UFColorBasemapStandard/VectorTileServer",
      {
        pane: "campusPane",
        maxZoom: 20,
      },
    );

    base.addTo(map);
    ufcampus.addTo(map);
    polygon.addTo(map);
    window.L.marker([29.64641, -82.34788], {
      pane: "markerPane",
    })
      .bindPopup("Address: 655 Reitz Union Dr Gainesville, FL 32611")
      .addTo(map);

    return () => {
      if (mapRef.current) {
        map.remove();
      }
    };
  }, []);

  return (
    <>
      <HeadContent />
      <div className="relative flex flex-col space-y-5 max-w-3xl mx-auto pb-10">
        <div className="relative flex gap-1">
          <span className="absolute text-xl font-bold text-text-main mb-1 top-15">
            Welcome, {user.name}!
          </span>

          <div className="absolute -top-5 right-0">
            <img className="w-50 select-none" src="/assets/XII_Sticker.png" />
          </div>
        </div>

        <div className="space-y-3 mt-25">
          <p>
            SwampHacks brings together 300+ students each year for 36 hours of
            creativity, collaboration, and innovation. Recognized for excellence
            by UF’s Herbert Wertheim College of Engineering, SwampHacks offers
            hands-on workshops, mentorship, and community-building activities
            that help hackers grow their skills and bring their ideas to life.
            Whether you’re a first-time hacker or a seasoned coder, SwampHacks
            is the place to build, connect, and inspire.
          </p>

          <div className="space-y-2">
            <div className="flex items-center gap-1">
              <TablerCalendarDue />
              <p>October 16 - 18, 2026</p>
            </div>
          </div>
        </div>

        <div className="space-y-2">
          <div>
            <p className="text-xl">📅 SwampHacks XII Timeline</p>
            <p className="text-text-secondary">Deadlines may change</p>
          </div>
          <div className="mt-5">
            <Timeline />
          </div>
        </div>

        <div className="space-y-2">
          <p className="text-xl">📍Event Location</p>
          <p className="my-2">
            Address: 655 Reitz Union Dr Gainesville, FL 32611
          </p>
          <div ref={mapRef} className="block h-100 rounded-md" id="map"></div>
          <p className="text-text-secondary">University of Florida Campus</p>
        </div>

        <OnboardingModal
          isOpen={isModalOpen}
          onOpenChange={setIsModalOpen}
        ></OnboardingModal>
      </div>
    </>
  );
}
