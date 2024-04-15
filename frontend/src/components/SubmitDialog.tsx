// SubmitDialog.tsx
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogDescription, DialogFooter } from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";

interface SubmitDialogProps {
    open: boolean;
    onClose: () => void;
    resultImage: string | null;
}

export function SubmitDialog({ open, onClose, resultImage }: SubmitDialogProps) {
    return (
        <Dialog open={open} onOpenChange={onClose}>
            <DialogContent>
                <DialogHeader>
                    <DialogTitle>Swap Result</DialogTitle>
                    <DialogDescription>
                        Here is the result of your image swap:
                    </DialogDescription>
                </DialogHeader>
                {resultImage ? (
                    <img
                        alt="Swap Result"
                        className="object-cover rounded-lg mb-4"
                        height={200}
                        src={resultImage}
                        style={{
                            aspectRatio: "200/200",
                            objectFit: "cover",
                        }}
                        width={400}
                    />
                ) : (
                    <div className="flex justify-center items-center h-48">
                        <div className="animate-pulse rounded-full bg-gray-400 h-12 w-12"></div>
                        <h1 className="ml-2">Loading...</h1>
                    </div>
                )}
            </DialogContent>
            <DialogFooter>
                <Button onClick={onClose}>Close</Button>
            </DialogFooter>
        </Dialog>
    );
}
