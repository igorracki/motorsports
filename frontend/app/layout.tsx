import React from "react"
import type { Metadata } from 'next'
import { Geist, Geist_Mono } from 'next/font/google'
import './globals.css'
import { AuthProvider } from "@/hooks/useAuth"
import { SettingsProvider } from "@/hooks/useSettings"
import { ApiProvider } from "@/components/providers/api-provider"
import { ConfigProvider } from "@/components/providers/config-provider"
import { ModalProvider } from "@/components/providers/modal-provider"

const geistSans = Geist({
  variable: "--font-geist-sans",
  subsets: ["latin"],
});

const geistMono = Geist_Mono({
  variable: "--font-geist-mono",
  subsets: ["latin"],
});

export const metadata: Metadata = {
  title: 'F1 Weekends',
}

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode
}>) {
  return (
    <html lang="en" className="dark">
      <body className={`${geistSans.variable} ${geistMono.variable} font-sans antialiased`}>
        <ApiProvider>
          <ConfigProvider>
            <AuthProvider>
              <SettingsProvider>
                <ModalProvider>
                  {children}
                </ModalProvider>
              </SettingsProvider>
            </AuthProvider>
          </ConfigProvider>
        </ApiProvider>
      </body>
    </html>
  )
}
