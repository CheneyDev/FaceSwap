"use client";
import React, { useRef, useState } from 'react';
import { AvatarImage, AvatarFallback, Avatar } from "@/components/ui/avatar";
import { Button } from "@/components/ui/button";

export function SwapPage() {
    const [image1, setImage1] = useState<string | null>(null);
    const [image2, setImage2] = useState<string | null>(null);
    const fileInputRef1 = useRef<HTMLInputElement | null>(null);
    const fileInputRef2 = useRef<HTMLInputElement | null>(null);

    const handleImageUpload = (event: React.ChangeEvent<HTMLInputElement>, setImage: (value: string | null) => void) => {
        const file = event.target.files?.[0];
        const reader = new FileReader();

        reader.onload = () => {
            setImage(reader.result as string);
        };

        if (file) {
            reader.readAsDataURL(file);
        }
    };

    const handleDivClick = (ref: React.RefObject<HTMLInputElement>) => {
        ref.current?.click();
    };

    const handleSubmit = async () => {
        if (image1 && image2 && fileInputRef1.current && fileInputRef2.current) {
            const formData = new FormData();
            // @ts-ignore
            formData.append('image1', fileInputRef1.current.files[0]);
            // @ts-ignore
            formData.append('image2', fileInputRef2.current.files[0]);

            try {
                const response = await fetch('/api/upload', {
                    method: 'POST',
                    body: formData,
                    credentials: 'include', // 发送 cookie
                });

                if (response.ok) {
                    console.log('Images uploaded successfully');
                    // 处理上传成功的逻辑
                } else {
                    console.error('Image upload failed');
                    // 处理上传失败的逻辑
                }
            } catch (error) {
                console.error('Error uploading images:', error);
                // 处理错误
            }
        }
    };

    return (
        <div className="max-w-4xl mx-auto p-4">
            <header className="flex justify-between items-center border-b pb-2">
                <h1 className="text-xl font-bold">
          <span className="bg-gradient-to-r from-purple-400 to-pink-600 bg-clip-text text-transparent">
            FaceSwap AI
          </span>
                </h1>
                <div className="flex items-center space-x-2">
                    <span>username</span>
                    <Avatar>
                        <AvatarImage alt="user avatar" src="/placeholder.svg?height=32&width=32" />
                        <AvatarFallback>U</AvatarFallback>
                    </Avatar>
                </div>
            </header>
            <main className="mt-8 h-screen">
                <div className="grid grid-cols-2 gap-4 h-fit">
                    <div
                        className="border-solid border border-gray-400 rounded-lg p-4 h-fit flex flex-col items-center text-sm cursor-pointer"
                        onClick={() => handleDivClick(fileInputRef1)}
                        role="button"
                        tabIndex={0}
                    >
                        <img
                            alt="Uploaded Image 1"
                            className="object-cover rounded-lg mb-4"
                            height={200}
                            src={image1 || "/placeholder.svg"}
                            style={{
                                aspectRatio: "200/200",
                                objectFit: "cover",
                            }}
                            width={400}
                        />
                        <input
                            type="file"
                            ref={fileInputRef1}
                            onChange={(e) => handleImageUpload(e, setImage1)}
                            style={{ display: 'none' }}
                        />
                        Click to upload image
                    </div>
                    <div
                        className="border-solid border border-gray-400 rounded-lg p-4 h-fit flex flex-col items-center text-sm cursor-pointer"
                        onClick={() => handleDivClick(fileInputRef2)}
                        role="button"
                        tabIndex={0}
                    >
                        <img
                            alt="Uploaded Image 2"
                            className="object-cover rounded-lg mb-4"
                            height={200}
                            src={image2 || "/placeholder.svg"}
                            style={{
                                aspectRatio: "200/200",
                                objectFit: "cover",
                            }}
                            width={400}
                        />
                        <input
                            type="file"
                            ref={fileInputRef2}
                            onChange={(e) => handleImageUpload(e, setImage2)}
                            style={{ display: 'none' }}
                        />
                        Click to upload image
                    </div>
                </div>
                <div className="flex justify-center mt-8">
                    <Button className="w-32" onClick={handleSubmit}>
                        Submit
                    </Button>
                </div>
            </main>
        </div>
    );
}