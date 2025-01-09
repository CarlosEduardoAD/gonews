import { Button } from "@/components/button";
import Image from "next/image";
import Link from "next/link";
import { redirect, useSearchParams } from "next/navigation";
import ResendButton from "./ResendButton";

export default async function EmailConfirmation({
  searchParams,
}: {
  searchParams: Promise<{ token: string }>;
}) {
  const token = (await searchParams).token;

  if (!token) {
    redirect("/");
  }

  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-50">
      <div className="max-w-md w-full space-y-8 p-8 bg-white rounded-lg shadow-lg">
        <div className="text-center">
          <div className="mx-auto w-24 h-24 relative mb-4">
            <Image
              src="/images/gopher_speed_of_light.png" // Você precisará adicionar esta imagem
              alt="Email enviado"
              fill
              className="object-contain"
            />
          </div>

          <h2 className="mt-6 text-3xl font-bold text-gray-900">
            Verifique seu e-mail
          </h2>

          <p className="mt-2 text-sm text-gray-600">
            Enviamos um link de confirmação para o seu e-mail. Por favor,
            verifique sua caixa de entrada e clique no link para confirmar seu
            endereço.
          </p>

          <div className="mt-6">
            <p className="text-xs text-gray-500">
              Não recebeu o e-mail?{" "}
              <ResendButton token={token} />
            </p>
          </div>

          <div className="mt-8">
            <Link
              href="/"
              className="text-sm font-medium text-indigo-600 hover:text-indigo-500"
            >
              ← Voltar para a página inicial
            </Link>
          </div>
        </div>
      </div>
    </div>
  );
}

