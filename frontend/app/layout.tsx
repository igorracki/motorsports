import React, { Suspense } from "react"
import type { Metadata } from 'next'
import { Geist, Geist_Mono } from 'next/font/google'
import './globals.css'
import { SeasonProvider } from "@/hooks/SeasonContext"

const _geist = Geist({ subsets: ["latin"] });
const _geistMono = Geist_Mono({ subsets: ["latin"] });

export const metadata: Metadata = {
  title: 'F1 Data Hub',
  description: 'Formula 1 Calendar, Results and Predictions',
}

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode
}>) {
  return (
    <html lang="en" className="dark">
      <body className={`font-sans antialiased`}>
        <Suspense fallback={null}>
          <SeasonProvider>
            {children}
          </SeasonProvider>
        </Suspense>
      </body>
    </html>
  )
}
