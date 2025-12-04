import { Header } from "@/components/layout/Header"
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
      <Header />
      <div className="my-2 mx-1">
        <HeroSection />
      </div>

      <div className="container mx-auto px-4 md:px-6">
        <DetailedFeatures />
        <ObjectionsSection />
        <SocialProofSection />
        <FAQSection />
        <PricingSection />
        <FinalCTASection />
      </div>
    </div>
  )
}

export default Home
