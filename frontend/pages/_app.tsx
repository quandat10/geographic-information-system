import "../styles/globals.css";
import type { AppProps } from "next/app";
import { MapContextProvider } from "../context/context";

function MyApp({ Component, pageProps }: AppProps) {
  return <MapContextProvider>
    <Component {...pageProps} />
  </MapContextProvider>;
}
export default MyApp;
