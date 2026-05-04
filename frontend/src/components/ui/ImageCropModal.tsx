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

type DragMode = 'move' | 'resize-tl' | 'resize-tr' | 'resize-bl' | 'resize-br' | 'resize-t' | 'resize-b' | 'resize-l' | 'resize-r' | null;

export function ImageCropModal({ image, onSave, onClose, title = 'Обрезать изображение' }: ImageCropModalProps) {
  const canvasRef = useRef<HTMLCanvasElement>(null);
  const [loadedImage, setLoadedImage] = useState<HTMLImageElement | null>(null);
  const [isDragging, setIsDragging] = useState(false);
  const [dragMode, setDragMode] = useState<DragMode>(null);
  const [dragStartPos, setDragStartPos] = useState({ x: 0, y: 0 });
  const [initialCropArea, setInitialCropArea] = useState({ x: 0, y: 0, width: 0, height: 0 });
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

  // Определяет режим перетаскивания на основе позиции клика
  const getDragMode = (x: number, y: number): DragMode => {
    const cornerSize = 30; // Размер зоны для захвата угла
    const edgeSize = 20; // Размер зоны для захвата стороны
    const { x: cropX, y: cropY, width, height } = cropArea;

    const nearLeft = Math.abs(x - cropX) < cornerSize;
    const nearRight = Math.abs(x - (cropX + width)) < cornerSize;
    const nearTop = Math.abs(y - cropY) < cornerSize;
    const nearBottom = Math.abs(y - (cropY + height)) < cornerSize;

    const onLeft = Math.abs(x - cropX) < edgeSize;
    const onRight = Math.abs(x - (cropX + width)) < edgeSize;
    const onTop = Math.abs(y - cropY) < edgeSize;
    const onBottom = Math.abs(y - (cropY + height)) < edgeSize;

    // Проверяем углы (приоритет выше чем стороны)
    if (nearLeft && nearTop) return 'resize-tl';
    if (nearRight && nearTop) return 'resize-tr';
    if (nearLeft && nearBottom) return 'resize-bl';
    if (nearRight && nearBottom) return 'resize-br';

    // Проверяем стороны (только если НЕ в зонах углов)
    if (onTop && !nearLeft && !nearRight) return 'resize-t';
    if (onBottom && !nearLeft && !nearRight) return 'resize-b';
    if (onLeft && !nearTop && !nearBottom) return 'resize-l';
    if (onRight && !nearTop && !nearBottom) return 'resize-r';

    // Проверяем центр (перемещение)
    if (x > cropX + edgeSize && x < cropX + width - edgeSize &&
        y > cropY + edgeSize && y < cropY + height - edgeSize) {
      return 'move';
    }

    return null;
  };

  const handleTouchStart = (e: React.TouchEvent) => {
    if (!canvasRef.current) return;

    const canvas = canvasRef.current;
    const rect = canvas.getBoundingClientRect();
    const touch = e.touches[0];
    const x = touch.clientX - rect.left;
    const y = touch.clientY - rect.top;

    const mode = getDragMode(x, y);
    if (mode) {
      setIsDragging(true);
      setDragMode(mode);
      setDragStartPos({ x, y });
      setInitialCropArea({ ...cropArea });
    }
  };

  const handleTouchMove = (e: React.TouchEvent) => {
    if (!isDragging || !canvasRef.current || !dragMode) return;

    const canvas = canvasRef.current;
    const rect = canvas.getBoundingClientRect();
    const touch = e.touches[0];
    const x = touch.clientX - rect.left;
    const y = touch.clientY - rect.top;

    updateCropArea(x, y);
  };

  const handleTouchEnd = () => {
    setIsDragging(false);
    setDragMode(null);
  };

  // Mouse events для поддержки на десктопе
  const handleMouseDown = (e: React.MouseEvent) => {
    if (!canvasRef.current) return;

    const canvas = canvasRef.current;
    const rect = canvas.getBoundingClientRect();
    const x = e.clientX - rect.left;
    const y = e.clientY - rect.top;

    const mode = getDragMode(x, y);
    if (mode) {
      setIsDragging(true);
      setDragMode(mode);
      setDragStartPos({ x, y });
      setInitialCropArea({ ...cropArea });
    }
  };

  const handleMouseMove = (e: React.MouseEvent) => {
    if (!canvasRef.current) return;

    const canvas = canvasRef.current;
    const rect = canvas.getBoundingClientRect();
    const x = e.clientX - rect.left;
    const y = e.clientY - rect.top;

    // Изменяем курсор в зависимости от позиции
    if (!isDragging) {
      const mode = getDragMode(x, y);
      canvas.style.cursor = getCursorForMode(mode);
    }

    if (!isDragging || !dragMode) return;
    updateCropArea(x, y);
  };

  const handleMouseUp = () => {
    setIsDragging(false);
    setDragMode(null);
  };

  const handleMouseLeave = () => {
    setIsDragging(false);
    setDragMode(null);
  };

  // Получает курсор для режима
  const getCursorForMode = (mode: DragMode): string => {
    switch (mode) {
      case 'move': return 'move';
      case 'resize-tl': case 'resize-br': return 'nwse-resize';
      case 'resize-tr': case 'resize-bl': return 'nesw-resize';
      case 'resize-t': case 'resize-b': return 'ns-resize';
      case 'resize-l': case 'resize-r': return 'ew-resize';
      default: return 'default';
    }
  };

  // Обновляет область обрезки в зависимости от режима
  const updateCropArea = (x: number, y: number) => {
    if (!canvasRef.current) return;

    const canvas = canvasRef.current;
    const dx = x - dragStartPos.x;
    const dy = y - dragStartPos.y;
    const minSize = 100;

    let newCrop = { ...initialCropArea };

    switch (dragMode) {
      case 'move':
        newCrop.x = Math.max(0, Math.min(canvas.width - initialCropArea.width, initialCropArea.x + dx));
        newCrop.y = Math.max(0, Math.min(canvas.height - initialCropArea.height, initialCropArea.y + dy));
        break;

      case 'resize-tl':
        newCrop.x = Math.min(initialCropArea.x + dx, initialCropArea.x + initialCropArea.width - minSize);
        newCrop.y = Math.min(initialCropArea.y + dy, initialCropArea.y + initialCropArea.height - minSize);
        newCrop.width = initialCropArea.width - (newCrop.x - initialCropArea.x);
        newCrop.height = initialCropArea.height - (newCrop.y - initialCropArea.y);
        newCrop.x = Math.max(0, newCrop.x);
        newCrop.y = Math.max(0, newCrop.y);
        break;

      case 'resize-tr':
        newCrop.y = Math.min(initialCropArea.y + dy, initialCropArea.y + initialCropArea.height - minSize);
        newCrop.width = Math.max(minSize, Math.min(initialCropArea.width + dx, canvas.width - initialCropArea.x));
        newCrop.height = initialCropArea.height - (newCrop.y - initialCropArea.y);
        newCrop.y = Math.max(0, newCrop.y);
        break;

      case 'resize-bl':
        newCrop.x = Math.min(initialCropArea.x + dx, initialCropArea.x + initialCropArea.width - minSize);
        newCrop.width = initialCropArea.width - (newCrop.x - initialCropArea.x);
        newCrop.height = Math.max(minSize, Math.min(initialCropArea.height + dy, canvas.height - initialCropArea.y));
        newCrop.x = Math.max(0, newCrop.x);
        break;

      case 'resize-br':
        newCrop.width = Math.max(minSize, Math.min(initialCropArea.width + dx, canvas.width - initialCropArea.x));
        newCrop.height = Math.max(minSize, Math.min(initialCropArea.height + dy, canvas.height - initialCropArea.y));
        break;

      case 'resize-t':
        newCrop.y = Math.min(initialCropArea.y + dy, initialCropArea.y + initialCropArea.height - minSize);
        newCrop.height = initialCropArea.height - (newCrop.y - initialCropArea.y);
        newCrop.y = Math.max(0, newCrop.y);
        break;

      case 'resize-b':
        newCrop.height = Math.max(minSize, Math.min(initialCropArea.height + dy, canvas.height - initialCropArea.y));
        break;

      case 'resize-l':
        newCrop.x = Math.min(initialCropArea.x + dx, initialCropArea.x + initialCropArea.width - minSize);
        newCrop.width = initialCropArea.width - (newCrop.x - initialCropArea.x);
        newCrop.x = Math.max(0, newCrop.x);
        break;

      case 'resize-r':
        newCrop.width = Math.max(minSize, Math.min(initialCropArea.width + dx, canvas.width - initialCropArea.x));
        break;
    }

    setCropArea(newCrop);
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
            onMouseDown={handleMouseDown}
            onMouseMove={handleMouseMove}
            onMouseUp={handleMouseUp}
            onMouseLeave={handleMouseLeave}
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
