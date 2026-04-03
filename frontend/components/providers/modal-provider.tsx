"use client";

import { createContext, useContext, useState, ReactNode } from "react";
import { LoginModal } from "@/components/features/auth/LoginModal";

interface ModalContextType {
  openLoginModal: (onSuccess?: () => void) => void;
  closeLoginModal: () => void;
}

const ModalContext = createContext<ModalContextType | undefined>(undefined);

export function ModalProvider({ children }: { children: ReactNode }) {
  const [isLoginOpen, setIsLoginOpen] = useState(false);
  const [onSuccessCallback, setOnSuccessCallback] = useState<(() => void) | undefined>();

  const openLoginModal = (onSuccess?: () => void) => {
    setOnSuccessCallback(() => onSuccess);
    setIsLoginOpen(true);
  };

  const closeLoginModal = () => {
    setIsLoginOpen(false);
    // Clear callback after a short delay to allow transition to finish
    setTimeout(() => setOnSuccessCallback(undefined), 300);
  };

  const handleSuccess = () => {
    if (onSuccessCallback) {
      onSuccessCallback();
    }
  };

  return (
    <ModalContext.Provider value={{ openLoginModal, closeLoginModal }}>
      {children}
      <LoginModal
        isOpen={isLoginOpen}
        onClose={closeLoginModal}
        onSuccess={handleSuccess}
      />
    </ModalContext.Provider>
  );
}

export function useModal() {
  const context = useContext(ModalContext);
  if (context === undefined) {
    throw new Error("useModal must be used within a ModalProvider");
  }
  return context;
}
