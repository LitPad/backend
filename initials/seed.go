package initials

import (
	"github.com/LitPad/backend/models"
	"github.com/LitPad/backend/models/choices"
)

var (
	TAGS []string = []string{
		"Campus", "Revenge", "Second Chance", "Age gap", "Regret", "Misunderstanding",
		"Hatred", "Kiss", "Runaway", "Hate to love", "Strong female lead", "School",
		"Forgiveness", "Royal", "Arranged marriage", "Betrayal", "Secret crush", "Pregnant",
		"Character growth", "He", "She", "Destiny", "Painful love", "Rejected", "Enemies",
		"Prophecy", "Shifter", "Omega verse", "Alternate universe", "Controlling", "World domination",
		"Karma", "Beautiful female lead", "Hardship", "Hard-working protagonist", "Future",	
	}

	GENRES []string = []string{
		"Werewolf", "Romance", "Billonaire", "Fantasy", "Mafia",
		"Historical", "YA/Teen", "Paranormal", "Urban", "Sci-fi", "Chicklit",
	}

	SECTIONS []string = []string{"New", "Romance", "Paranormal", "Dark Romance", "Young Adult", "Long Reads"}

	SUBSECTIONS []string = []string{
		"Steamy Read", "Enemy to Lovers", "Editor's Pick", "Rising Stars", "Steamy & Wholesome",
		"Forbidden Love", "Bestselling Romances", "Werewolf Tales", "Vampire Chronicles", "Fated Mates", "Shapeshifter & Beasts",
		"Paranormal Romance", "Dark Magic", "Twisted Love", "Billonaire Dark Romances", "Dangerous Affair", "Revenge & Redemption",
		"Possessive Lovers", "Teen Drama", "Coming Of Age", "Young love", "Fantasy Adventures", "Epic Adventures",
		"Series You Can't Put Down", "Binge-Worthy Romances", "Complete Sagas", "Fantasy Epics", "Ultimate Page-Turners",
	}

	GIFTNAMES []string = []string{"Red rose", "Black dahlia", "Scroll", "Magic wand", "Wolf", "Baby Dragon"}

	BOOK = models.Book{
		Title:  "The Mysterious Island",
		Blurb:  "An adventure novel about survival and discovery on an uncharted island.",
		AgeDiscretion: choices.ATYPE_EIGHTEEN,
		CoverImage: "https://example.com/cover.jpg",
		Completed:  false,
	}

	PARAGRAPHS = []string{
		"The sun hung low over the horizon, casting golden streaks across the sky as Ethan, Mira, and Jackson climbed the steep rock face. It had been three days since their arrival on the island—a storm-tossed fate that left them stranded on this mysterious land with no sign of rescue. The island was lush, teeming with strange plants and unseen creatures that whispered through the dense jungle. Yet, it was the beacon of light they had spotted from the shore the night before that drove them toward the cliffs.",
		"Their journey had been fraught with danger and discovery, from the moment they washed ashore to the eerie silence of the forest. Ethan, the eldest of the three, had taken charge, his keen eyes scanning the horizon for any sign of life. Mira, the youngest, clung to his side, her eyes wide with wonder and fear. Jackson, the middle child, had been the first to spot the light, a flickering beacon that promised hope in the darkness.",
		"As they reached the summit, the sun dipped below the waves, casting long shadows across the rocky plateau. Ethan raised his hand to shield his eyes, squinting against the dying light. Mira gasped, her hand tightening on his arm as she pointed toward the source of the light. Jackson let out a whoop of joy, his voice echoing across the cliffs as he raced toward the edge.",
		"Below them, nestled in a hidden cove, lay the wreck of a ship—a massive hulk of wood and iron that had been torn asunder by the storm. The light they had seen was a lantern, swinging from the mast, its glow beckoning them down to the shore. Ethan felt a surge of hope as he gazed upon the wreck, his heart pounding in his chest. Mira clung to his side, her eyes wide with wonder as she took in the sight.",
		"Without a word, the three of them began the treacherous descent, picking their way down the cliff face toward the cove. The path was narrow and steep, the rocks slick with spray from the crashing waves below. Ethan led the way, his hands and feet finding purchase on the rough stone as he guided his siblings down the slope. Mira followed close behind, her breath coming in short gasps as she fought to keep up. Jackson brought up the rear, his eyes fixed on the lantern below as he scrambled after them.",
		"At last, they reached the bottom, their feet sinking into the soft sand of the cove. The wreck loomed before them, its shattered hull rising from the water like a ghostly specter. Ethan felt a shiver run down his spine as he gazed upon the wreck, his mind racing with questions. Who had been aboard the ship? What had happened to them? And most importantly, was there any hope of rescue?",
		"Mira tugged at his sleeve, her eyes wide with fear as she stared up at the wreck. 'What do we do now?' she whispered, her voice barely audible over the crash of the waves. Ethan turned to her, his face grim as he considered their options. 'We need to find a way inside,' he said, his voice firm. 'There may be survivors, or supplies we can use.' Mira nodded, her eyes shining with determination as she followed him toward the wreck.",
		"Jackson raced ahead, his feet pounding on the sand as he reached the wreck first. He scrambled up the side, his hands finding purchase on the splintered wood as he pulled himself over the rail. Ethan and Mira followed close behind, their hearts pounding in their chests as they climbed aboard. The deck was littered with debris, the remains of the storm that had torn the ship apart. Ethan scanned the wreckage, his eyes searching for any sign of life.",
		"As they made their way belowdecks, the air grew thick with the smell of salt and decay. The lantern cast a flickering light on the walls, revealing the shattered remains of the ship's interior. Ethan led the way, his eyes scanning the shadows for any sign of movement. Mira clung to his side, her breath coming in short gasps as she fought to keep up. Jackson brought up the rear, his eyes wide with wonder as he took in the sight.",
		"They searched every corner of the wreck, from the bow to the stern, but found no sign of life. The ship was deserted, its crew lost to the storm that had claimed it. Ethan felt a surge of despair as he gazed upon the empty decks, his heart heavy with the knowledge of their fate. Mira clung to his side, her eyes wide with fear as she took in the sight. Jackson stood beside them, his face grim as he surveyed the wreckage.",
		"At last, they reached the captain's cabin, the door hanging open on its hinges. Ethan pushed it aside, his heart pounding in his chest as he stepped inside. The room was dark, the air thick with the smell of salt and decay. Mira clung to his side, her eyes wide with fear as she scanned the shadows. Jackson stood beside them, his hand on Ethan's shoulder as he peered into the gloom.",
	}
)


