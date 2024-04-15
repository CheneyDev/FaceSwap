"use client";
import React, { useRef, useState } from 'react';
import { AvatarImage, AvatarFallback, Avatar } from "@/components/ui/avatar";
import { Button } from "@/components/ui/button";

export function SwapPage() {
    const [image1, setImage1] = useState(null);
    const [image2, setImage2] = useState(null);
    const fileInputRef1 = useRef(null);
    const fileInputRef2 = useRef(null);

    const handleImageUpload = (event, setImage) => {
        const file = event.target.files[0];
        const reader = new FileReader();

        reader.onload = () => {
            setImage(reader.result);
        };

        if (file) {
            reader.readAsDataURL(file);
        }
    };

    const handleDivClick = (ref) => {
        ref.current.click();
    };

    return (
        <div className="max-w-4xl mx-auto p-4">
            <header className="flex justify-between items-center border-b pb-2">
                <h1 className="text-xl font-bold">FaceSwap AI</h1>
                <div className="flex items-center space-x-2">
                    <span>username</span>
                    <Avatar>
                        <AvatarImage alt="user avatar" src="/placeholder.svg?height=32&width=32" />
                        <AvatarFallback>U</AvatarFallback>
                    </Avatar>
                </div>
            </header>
            <main className="mt-4 h-screen">
                <div className="grid grid-cols-2 gap-4 h-fit">
                    <div
                        className="border-dashed border-2 border-gray-400 rounded-lg p-4 h-fit flex flex-col items-center text-sm cursor-pointer"
                        onClick={() => handleDivClick(fileInputRef1)}
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
                        className="border-dashed border-2 border-gray-400 rounded-lg p-4 h-fit flex flex-col items-center text-sm cursor-pointer"
                        onClick={() => handleDivClick(fileInputRef2)}
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
                <div className="flex justify-center mt-4">
                    <Button className="w-32">Submit</Button>
                </div>
            </main>
        </div>
    );
}