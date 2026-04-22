// src/pages/Friends/components/ReferralLinkSection.tsx
import { useState } from 'react';
import { Card } from '@/components/ui/Card';
import { Button } from '@/components/ui/Button';
import bridge from '@/lib/platform/bridge';
import styles from './ReferralLinkSection.module.css';

interface ReferralLinkSectionProps {
  referralCode: string;
  referralLink: string;
  onCopy: () => void;
  onShare: (channel: string) => void;
}

export function ReferralLinkSection({
  referralCode,
  referralLink,
  onCopy,
  onShare,
}: ReferralLinkSectionProps) {
  const [copied, setCopied] = useState(false);

  const handleCopy = async () => {
    try {
      await navigator.clipboard.writeText(referralLink);
      setCopied(true);
      onCopy();
      setTimeout(() => setCopied(false), 2000);
    } catch (error) {
      console.error('[ReferralLinkSection] Failed to copy:', error);
    }
  };

  const handleShare = async () => {
    try {
      await bridge.send('VKWebAppShare', {
        link: referralLink,
      });
      onShare('vk');
    } catch (error) {
      console.error('[ReferralLinkSection] Failed to share:', error);
    }
  };

  return (
    <Card className={styles.card}>
      <h3 className={styles.title}>Твоя реферальная ссылка</h3>
      <p className={styles.description}>
        Поделись ссылкой с друзьями и получай награды
      </p>

      <div className={styles.codeContainer}>
        <div className={styles.code}>{referralCode}</div>
      </div>

      <div className={styles.linkContainer}>
        <input
          type="text"
          value={referralLink}
          readOnly
          className={styles.linkInput}
        />
      </div>

      <div className={styles.buttons}>
        <Button
          mode={copied ? 'secondary' : 'primary'}
          size="l"
          stretched
          onClick={handleCopy}
        >
          {copied ? '✓ Скопировано' : 'Скопировать ссылку'}
        </Button>

        <Button
          mode="secondary"
          size="l"
          stretched
          onClick={handleShare}
        >
          Поделиться
        </Button>
      </div>
    </Card>
  );
}
