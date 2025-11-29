import { HeroSection } from "@/components/home/HeroSection"
import { DetailedFeatures } from "@/components/home/DetailedFeatures"
import { ObjectionsSection } from "@/components/home/ObjectionsSection"
import { SocialProofSection } from "@/components/home/SocialProofSection"
import { FAQSection } from "@/components/home/FAQSection"
import { PricingSection } from "@/components/home/PricingSection"
import { FinalCTASection } from "@/components/home/FinalCTASection"

const Home = () => {
  return (
    <div className="flex flex-col">
      {/* Hero Section */}
      <div className="my-2 mx-1">
        <HeroSection />
      </div>

      {/* Main Content */}
      <div className="container mx-auto px-4 md:px-6">
        {/* Features */}
        <DetailedFeatures />

        {/* Objections */}
        <ObjectionsSection />

        {/* Social Proof */}
        <SocialProofSection />

        {/* FAQ */}
        <FAQSection />

        {/* Pricing */}
        <PricingSection />

        {/* Final CTA */}
        <FinalCTASection />
      </div>
    </div>
  )
}

export default Home
