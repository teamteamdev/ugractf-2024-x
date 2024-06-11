import { redirect } from "next/navigation";

export default function Jobs() {
  redirect("https://hh.ru/vacancy/" + VACANCY_ID);
}
