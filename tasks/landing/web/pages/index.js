import Head from "next/head";
import Link from "next/link";
import { Lilita_One, Noto_Sans } from "next/font/google";
import styles from "../styles/Home.module.css";

const lilita = Lilita_One({
  weight: "400",
  subsets: ["latin"],
  display: "swap",
});

const noto = Noto_Sans({
  subsets: ["cyrillic"],
  display: "swap",
});

export default function Home() {
  return (
    <main className={styles.main}>
      <Head>
        <title>Teemoorka.Network сейчас настраивается</title>
      </Head>

      <h1 className={`${styles.title} ${lilita.className}`}>Teemoorka.Network</h1>

      <h2 className={`${styles.inprogress} ${noto.className}`}>Сайт скоро откроется</h2>

      <Link className={`${styles.jobs} ${noto.className}`} href="/jobs">Приглашаем на работу</Link>
    </main>
  );
}
