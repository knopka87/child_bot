// src/components/ui/ImageCropModal.tsx
import { useState, useRef, useEffect } from 'react';
import { X, Check } from 'lucide-react';
import styles from './ImageCropModal.module.css';

interface ImageCropModalProps {
  image: string;
  onSave: (croppedImage: File) => void;
  onClose: () => void;
  title?: string;
}

export function ImageCropModal({ image, onSave, onClose, title = 'Обрезать изображение' }: ImageCropModalProps) {
  const canvasRef = useRef<HTMLCanvasElement>(null);
  const [loadedImage, setLoadedImage] = useState<HTMLImageElement | null>(null);
  const [isDragging, setIsDragging] = useState(false);
  const [cropArea, setCropArea] = useState({ x: 0, y: 0, width: 0, height: 0 });
  const [imageLoaded, setImageLoaded] = useState(false);

  useEffect(() => {
    const img = new Image();
    img.onload = () => {
      setLoadedImage(img);
      const canvas = canvasRef.current;
      if (!canvas) return;

      const ctx = canvas.getContext('2d');
      if (!ctx) return;

      // Масштабируем изображение чтобы поместилось в canvas
      const maxWidth = window.innerWidth - 40;
      const maxHeight = window.innerHeight - 200;
      let width = img.width;
      let height = img.height;

      if (width > maxWidth) {
        height = (height * maxWidth) / width;
        width = maxWidth;
      }

      if (height > maxHeight) {
        width = (width * maxHeight) / height;
        height = maxHeight;
      }

      canvas.width = width;
      canvas.height = height;

      // Рисуем изображение
      ctx.drawImage(img, 0, 0, width, height);

      // Устанавливаем начальную область обрезки (80% от размера)
      const margin = 0.1;
      setCropArea({
        x: width * margin,
        y: height * margin,
        width: width * (1 - 2 * margin),
        height: height * (1 - 2 * margin),
      });

      setImageLoaded(true);
    };
    img.src = image;
  }, [image]);

  useEffect(() => {
    if (!imageLoaded || !canvasRef.current || !loadedImage) return;

    const canvas = canvasRef.current;
    const ctx = canvas.getContext('2d');
    if (!ctx) return;

    // Очищаем canvas
    ctx.clearRect(0, 0, canvas.width, canvas.height);

    // Рисуем изображение
    ctx.drawImage(loadedImage, 0, 0, canvas.width, canvas.height);

    // Затемняем области вне кропа
    ctx.fillStyle = 'rgba(0, 0, 0, 0.5)';
    ctx.fillRect(0, 0, canvas.width, cropArea.y);
    ctx.fillRect(0, cropArea.y, cropArea.x, cropArea.height);
    ctx.fillRect(cropArea.x + cropArea.width, cropArea.y, canvas.width - cropArea.x - cropArea.width, cropArea.height);
    ctx.fillRect(0, cropArea.y + cropArea.height, canvas.width, canvas.height - cropArea.y - cropArea.height);

    // Рисуем рамку кропа
    ctx.strokeStyle = '#6C5CE7';
    ctx.lineWidth = 3;
    ctx.strokeRect(cropArea.x, cropArea.y, cropArea.width, cropArea.height);

    // Рисуем углы
    const cornerSize = 20;
    ctx.fillStyle = '#6C5CE7';
    // Верхний левый
    ctx.fillRect(cropArea.x - 3, cropArea.y - 3, cornerSize, 6);
    ctx.fillRect(cropArea.x - 3, cropArea.y - 3, 6, cornerSize);
    // Верхний правый
    ctx.fillRect(cropArea.x + cropArea.width - cornerSize + 3, cropArea.y - 3, cornerSize, 6);
    ctx.fillRect(cropArea.x + cropArea.width - 3, cropArea.y - 3, 6, cornerSize);
    // Нижний левый
    ctx.fillRect(cropArea.x - 3, cropArea.y + cropArea.height - 3, cornerSize, 6);
    ctx.fillRect(cropArea.x - 3, cropArea.y + cropArea.height - cornerSize + 3, 6, cornerSize);
    // Нижний правый
    ctx.fillRect(cropArea.x + cropArea.width - cornerSize + 3, cropArea.y + cropArea.height - 3, cornerSize, 6);
    ctx.fillRect(cropArea.x + cropArea.width - 3, cropArea.y + cropArea.height - cornerSize + 3, 6, cornerSize);
  }, [cropArea, imageLoaded]);

  const handleTouchStart = () => {
    setIsDragging(true);
  };

  const handleTouchMove = (e: React.TouchEvent) => {
    if (!isDragging || !canvasRef.current) return;

    const canvas = canvasRef.current;
    const rect = canvas.getBoundingClientRect();
    const touch = e.touches[0];
    const x = touch.clientX - rect.left;
    const y = touch.clientY - rect.top;

    // Обновляем размер области обрезки
    const newWidth = Math.max(100, Math.min(x - cropArea.x, canvas.width - cropArea.x));
    const newHeight = Math.max(100, Math.min(y - cropArea.y, canvas.height - cropArea.y));

    setCropArea(prev => ({
      ...prev,
      width: newWidth,
      height: newHeight,
    }));
  };

  const handleTouchEnd = () => {
    setIsDragging(false);
  };

  const handleSave = async () => {
    if (!canvasRef.current || !loadedImage) return;

    const canvas = canvasRef.current;
    const img = loadedImage;

    // Создаём новый canvas для обрезанного изображения
    const croppedCanvas = document.createElement('canvas');
    const ctx = croppedCanvas.getContext('2d');
    if (!ctx) return;

    // Вычисляем соотношение между оригинальным изображением и canvas
    const scaleX = img.width / canvas.width;
    const scaleY = img.height / canvas.height;

    croppedCanvas.width = cropArea.width * scaleX;
    croppedCanvas.height = cropArea.height * scaleY;

    // Рисуем обрезанную область
    ctx.drawImage(
      img,
      cropArea.x * scaleX,
      cropArea.y * scaleY,
      cropArea.width * scaleX,
      cropArea.height * scaleY,
      0,
      0,
      croppedCanvas.width,
      croppedCanvas.height
    );

    // Конвертируем в Blob и затем в File
    croppedCanvas.toBlob((blob) => {
      if (!blob) return;
      const file = new File([blob], 'cropped-image.jpg', { type: 'image/jpeg' });
      onSave(file);
    }, 'image/jpeg', 0.9);
  };

  return (
    <div className={styles.overlay}>
      <div className={styles.modal}>
        <div className={styles.header}>
          <h3 className={styles.title}>{title}</h3>
          <button onClick={onClose} className={styles.closeButton}>
            <X size={24} />
          </button>
        </div>

        <div className={styles.canvasContainer}>
          <canvas
            ref={canvasRef}
            className={styles.canvas}
            onTouchStart={handleTouchStart}
            onTouchMove={handleTouchMove}
            onTouchEnd={handleTouchEnd}
          />
        </div>

        <p className={styles.hint}>
          Растяните рамку чтобы выделить нужную область
        </p>

        <div className={styles.actions}>
          <button onClick={onClose} className={styles.cancelButton}>
            Отмена
          </button>
          <button onClick={handleSave} className={styles.saveButton}>
            <Check size={18} />
            Применить
          </button>
        </div>
      </div>
    </div>
  );
}
