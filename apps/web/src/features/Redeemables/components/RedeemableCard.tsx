import { Card } from "@/components/ui/Card";
import { DialogTrigger, Button } from "react-aria-components";
import { RedeemableDetailsModal } from "./RedeemableDetailsModal";

export interface RedeemableCardProps {
  id: string;
  name: string;
  totalStock: number;
  maxUserAmount: number;
  totalRedeemed: number;
  eventId: string;
}

export function RedeemableCard({
  id,
  name,
  totalStock,
  maxUserAmount,
  totalRedeemed,
  eventId,
}: RedeemableCardProps) {
  const remaining = totalStock - (totalRedeemed as number);
  const percentageRemaining =
    totalStock > 0 ? ((remaining / totalStock) * 100).toFixed(0) : 0;

  return (
    <Card className="border hover:shadow-md transition-shadow">
      <div className="p-6 w-full">
        <div className="flex justify-between items-start mb-4">
          <h3 className="text-xl font-semibold text-text-main">{name}</h3>
          <div className="text-sm text-text-secondary">
            {remaining} / {totalStock} remaining
          </div>
        </div>

        <div className="space-y-2">
          <div className="flex items-center justify-between text-sm">
            <span className="text-text-secondary">Max per user:</span>
            <span className="text-text-main font-medium">{maxUserAmount}</span>
          </div>

          <div className="w-full bg-neutral-200 dark:bg-neutral-700 rounded-full h-2">
            <div
              className="bg-violet-600 h-2 rounded-full transition-all"
              style={{ width: `${percentageRemaining}%` }}
            />
          </div>

          <div className="flex justify-end mt-4">
            <DialogTrigger>
              <Button className="px-4 py-2 bg-violet-600 text-white rounded hover:bg-violet-700 transition-colors font-medium cursor-pointer">
                Redeem
              </Button>
              <RedeemableDetailsModal
                id={id}
                name={name}
                totalStock={totalStock}
                maxUserAmount={maxUserAmount}
                totalRedeemed={totalRedeemed as number}
                eventId={eventId}
              />
            </DialogTrigger>
          </div>
        </div>
      </div>
    </Card>
  );
}
